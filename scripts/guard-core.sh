#!/usr/bin/env bash
# Guard decision core (spec: agent-boundary). Shared by the Claude Code adapter
# (.claude/hooks/guard.sh) and the Antigravity adapter (.agents/hooks/guard-agy.sh).
#
# Library:      source scripts/guard-core.sh
#               guard_decide TOOL PROJECT INPUT...
#                  → GUARD_DECISION (allow|ask|deny) and GUARD_REASON
#               guard_log TOOL DECISION REASON PROJECT
# Command line: scripts/guard-core.sh TOOL PROJECT INPUT...
#                  → "decision<TAB>reason"
#
# TOOL uses the Claude Code names: Bash, Read, Glob, Grep, Edit, Write, NotebookEdit.
# For Bash, INPUT is the command; for the others, one or more paths.
# "allow" only means the guard has no objection: the adapters print nothing and the final
# decision is left to the harness rules. The guard never approves anything on its own.
# GUARD_CWD (optional) is the directory relative paths are resolved from.
#
# Known limits: this is text analysis, not a security boundary. Whatever it cannot resolve
# (variables, substitutions, interpreters) becomes "ask", never "allow".
# See docs/harness/local-guardrails.md.

# The patterns below with literal quotes and backslashes are intentional.
# shellcheck disable=SC2016,SC1003,SC2088
_G_US=$'\x1f'

_g_deny() { [[ -n $_R_DENY ]] || _R_DENY=$1; }
_g_ask() { [[ -n $_R_ASK ]] || _R_ASK=$1; }

# ------------------------------------------------------------------ paths

# _g_secret_path PATH → 0 if the path is a secret (.env*, *.pem, *.key, secrets/, credentials)
_g_secret_path() {
  local p=${1,,} base home=${HOME,,}
  p=${p%/}
  base=${p##*/}
  case $base in .env* | *.pem | *.key) return 0 ;; esac
  case "/$p/" in */secrets/*) return 0 ;; esac
  if [[ -n $home ]]; then
    case $p in
      "$home"/.ssh | "$home"/.ssh/* | "$home"/.aws | "$home"/.aws/* | \
        "$home"/.config/gcloud | "$home"/.config/gcloud/* | "$home"/.kube | "$home"/.kube/* | \
        "$home"/.gemini | "$home"/.gemini/* | "$home"/.claude | "$home"/.claude/*) return 0 ;;
    esac
  fi
  return 1
}

# _g_protected REL → 0 if the path (relative to the project) is harness configuration or an
# agent instruction file
_g_protected() {
  case $1 in
    .claude | .claude/* | .agents | .agents/* | .github/workflows | .github/workflows/* | \
      Makefile | scripts/guard*.sh | scripts/harness-test.sh | AGENTS.md | CLAUDE.md | GEMINI.md) return 0 ;;
  esac
  return 1
}

# _g_scan_value VALUE MUT [FORCE] — records the value as a candidate for the path check.
# MUT=1 when the command may change the file. FORCE=1 treats any value as a path.
_g_scan_value() {
  local v=$1 mut=$2 force=${3:-0}
  [[ -n $v ]] || return 0
  case $v in
    "~") v=$HOME ;;
    "~/"*) v=$HOME/${v:2} ;;
    "~"*) _g_ask "~user not resolved"; return 0 ;;
  esac
  if _g_secret_path "$v"; then
    _g_deny "secret access (path masked)"
    return 0
  fi
  case $v in /dev/null | /dev/stdin | /dev/stdout | /dev/stderr) return 0 ;; esac
  if ((force)) || [[ $v == /* || $v == */* || $v == . || $v == .. ]] || [[ -L $_G_CWD/$v ]]; then
    _C_VAL+=("$v")
    _C_MUT+=("$mut")
  fi
}

# _g_scan_word WORD MUT — a word of a command (argument, assignment or redirect target)
_g_scan_word() {
  local w=$1 mut=$2
  [[ -n $w && $w != *[[:space:]]* ]] || return 0
  if [[ $w == -* ]]; then
    if [[ $w == *=* ]]; then
      _g_scan_value "${w#*=}" "$mut"
    elif [[ $w =~ ^-[A-Za-z]([/~].*)$ ]]; then
      _g_scan_value "${BASH_REMATCH[1]}" "$mut"
    fi
    return 0
  fi
  _g_scan_value "$w" "$mut" "$mut"
  if [[ $w == *=* ]]; then
    _g_scan_value "${w#*=}" "$mut"
  fi
}

# _g_scan_cmdword WORD — the word that names the command
_g_scan_cmdword() {
  case $1 in /usr/* | /bin/* | /sbin/* | /opt/*) return 0 ;; esac
  _g_scan_word "$1" 0
}

# _g_resolve — resolves all candidates at once (realpath -m: handles .., ~ and symbolic links)
_g_resolve() {
  ((${#_C_VAL[@]} > 0)) || return 0
  local -a abs=() res=()
  local v i r proj
  for v in "${_C_VAL[@]}"; do
    if [[ $v == /* ]]; then abs+=("$v"); else abs+=("$_G_CWD/$v"); fi
  done
  mapfile -d '' -t res < <(realpath -m -z -- "$_G_PROJ" "${abs[@]}" 2>/dev/null)
  if ((${#res[@]} != ${#abs[@]} + 1)); then
    _g_ask "could not resolve the paths"
    return 0
  fi
  proj=${res[0]}
  for i in "${!abs[@]}"; do
    r=${res[i + 1]}
    if _g_secret_path "$r"; then
      _g_deny "secret access (path masked)"
    elif [[ $r != "$proj" && $r != "$proj"/* ]]; then
      _g_ask "outside the repository: $r"
    elif [[ ${_C_MUT[i]} == 1 ]] && _g_protected "${r#"$proj"/}"; then
      _g_ask "change to harness configuration: ${r#"$proj"/}"
    fi
  done
}

# ------------------------------------------------------------------- lexer
# _g_lex TEXT → fills _SEGS (the words of each command, separated by US) and _SEPS (the operator
# that follows each command). Markers: \x01<op> redirect whose target is the next word;
# \x02<op> complete redirect (2>&1). Inside single quotes the text is literal; $( and backtick
# become "ask" even inside double quotes. Heredoc bodies are not read as commands.

_gl_endword() {
  if ((inw)); then
    seg+="$word$_G_US"
    if ((hd_expect)); then
      hd_delims+=("$word")
      hd_strips+=("$hd_strip")
      hd_expect=0
    fi
    word=""
    inw=0
  fi
}

_gl_endseg() {
  _gl_endword
  if [[ -n $seg ]]; then
    _SEGS+=("$seg")
    _SEPS+=("$1")
  fi
  seg=""
}

_gl_dollar() {
  local m name
  nx=${s:i+1:1}
  if [[ $nx == '(' ]]; then
    _g_ask "command substitution"
    word+='$('
    i=$((i + 1))
  elif [[ ${s:i} =~ ^\$(\{([A-Za-z_][A-Za-z0-9_]*)\}|([A-Za-z_][A-Za-z0-9_]*)) ]]; then
    m=${BASH_REMATCH[0]}
    name=${BASH_REMATCH[2]}${BASH_REMATCH[3]}
    case $name in
      HOME) word+=$HOME ;;
      PWD | CLAUDE_PROJECT_DIR) word+=$_G_PROJ ;;
      *) _g_ask "unresolved environment variable (\$$name)"; word+=$m ;;
    esac
    i=$((i + ${#m} - 1))
  else
    [[ $nx == '{' ]] && _g_ask "parameter expansion"
    word+='$'
  fi
}

_gl_redir() {
  local op=$c mk=$'\x01'
  nx=${s:i+1:1}
  if [[ $c == '&' ]]; then
    _gl_endword
    op='&>'
    i=$((i + 1))
    if [[ ${s:i+1:1} == '>' ]]; then
      op='&>>'
      i=$((i + 1))
    fi
  else
    if [[ $nx == '(' ]]; then
      _g_ask "process substitution"
      inw=1
      word+="$c("
      i=$((i + 1))
      return 0
    fi
    if ((inw)) && [[ $word =~ ^[0-9]+$ ]]; then
      word=""
      inw=0
    else
      _gl_endword
    fi
    if [[ $c == '>' ]]; then
      case $nx in
        '>') op='>>'; i=$((i + 1)) ;;
        '|') op='>|'; i=$((i + 1)) ;;
        '&') op='>&'; i=$((i + 1)) ;;
      esac
    else
      case $nx in
        '<')
          op='<<'
          i=$((i + 1))
          if [[ ${s:i+1:1} == '<' ]]; then
            op='<<<'
            i=$((i + 1))
          elif [[ ${s:i+1:1} == '-' ]]; then
            op='<<-'
            i=$((i + 1))
          fi
          ;;
        '>') op='<>'; i=$((i + 1)) ;;
        '&') op='<&'; i=$((i + 1)) ;;
      esac
    fi
    if [[ $op == *'&' ]]; then
      while [[ ${s:i+1:1} =~ ^[0-9-]$ ]]; do
        op+=${s:i+1:1}
        mk=$'\x02'
        i=$((i + 1))
      done
    fi
    if [[ $op == '<<' || $op == '<<-' ]]; then
      _g_ask "heredoc (content not analyzed)"
      hd_expect=1
      hd_strip=0
      [[ $op == '<<-' ]] && hd_strip=1
    fi
  fi
  seg+="$mk$op$_G_US"
}

# consumes the pending heredoc bodies, starting from the newline at s[i]
_gl_heredoc() {
  local rest=${s:i+1} line eof=0 k d strip
  for k in "${!hd_delims[@]}"; do
    d=${hd_delims[k]}
    strip=${hd_strips[k]}
    while :; do
      if [[ $rest == *$'\n'* ]]; then
        line=${rest%%$'\n'*}
        rest=${rest#*$'\n'}
      else
        line=$rest
        rest=""
        eof=1
      fi
      ((strip)) && line=${line#"${line%%[!$'\t']*}"}
      [[ $line == "$d" ]] && break
      ((eof)) && break
      _HD_ALL+="$line"$'\n'
    done
  done
  hd_delims=()
  hd_strips=()
  i=$((n - ${#rest} - 1))
}

_g_lex() {
  local s=$1 n=${#1} i=0 c nx st=n word="" inw=0 seg="" rest
  local hd_expect=0 hd_strip=0
  local -a hd_delims=() hd_strips=()
  _SEGS=()
  _SEPS=()
  while ((i < n)); do
    c=${s:i:1}
    case $st in
      s)
        if [[ $c == "'" ]]; then st=n; else word+=$c; fi
        ;;
      d)
        case $c in
          '"') st=n ;;
          '\')
            nx=${s:i+1:1}
            if [[ $nx == '$' || $nx == '`' || $nx == '"' || $nx == '\' ]]; then
              word+=$nx
              i=$((i + 1))
            elif [[ $nx == $'\n' ]]; then
              i=$((i + 1))
            else
              word+=$c
            fi
            ;;
          '$') _gl_dollar ;;
          '`') _g_ask "command substitution (backtick)"; word+=$c ;;
          *) word+=$c ;;
        esac
        ;;
      n)
        case $c in
          "'") st=s; inw=1 ;;
          '"') st=d; inw=1 ;;
          '\')
            nx=${s:i+1:1}
            if [[ $nx == $'\n' ]]; then
              i=$((i + 1))
            elif [[ -n $nx ]]; then
              word+=$nx
              inw=1
              i=$((i + 1))
            fi
            ;;
          ' ' | $'\t' | $'\r') _gl_endword ;;
          $'\n')
            _gl_endseg $'\n'
            if ((${#hd_delims[@]})); then _gl_heredoc; fi
            ;;
          ';' | '(' | ')') _gl_endseg ';' ;;
          '&')
            nx=${s:i+1:1}
            if [[ $nx == '&' ]]; then
              _gl_endseg '&&'
              i=$((i + 1))
            elif [[ $nx == '>' ]]; then
              _gl_redir
            else
              _gl_endseg '&'
            fi
            ;;
          '|')
            nx=${s:i+1:1}
            if [[ $nx == '|' ]]; then
              _gl_endseg '||'
              i=$((i + 1))
            elif [[ $nx == '&' ]]; then
              _gl_endseg '|'
              i=$((i + 1))
            else
              _gl_endseg '|'
            fi
            ;;
          '<' | '>') _gl_redir ;;
          '#')
            if ((inw)); then
              word+=$c
            else
              rest=${s:i}
              rest=${rest%%$'\n'*}
              i=$((i + ${#rest} - 1))
            fi
            ;;
          '$') inw=1; _gl_dollar ;;
          '`') inw=1; _g_ask "command substitution (backtick)"; word+=$c ;;
          *) word+=$c; inw=1 ;;
        esac
        ;;
    esac
    i=$((i + 1))
  done
  if [[ $st != n ]]; then
    _g_ask "unclosed quotes"
  fi
  _gl_endseg ""
}

# ------------------------------------------------------------------- rules

_g_rules_git() {
  local -a g=("$@")
  local gi=0 gs="" a
  while ((gi < ${#g[@]})); do
    a=${g[gi]}
    case $a in
      -C | -c | --git-dir | --work-tree | --namespace | --exec-path) gi=$((gi + 2)) ;;
      -*) gi=$((gi + 1)) ;;
      *) gs=$a; break ;;
    esac
  done
  local -a r=("${g[@]:gi+1}")
  case $gs in
    push)
      for a in "${r[@]}"; do
        case $a in
          --force* | --delete | --mirror | -d | +*) _g_deny "destructive git push"; return 0 ;;
          -[!-]*) if [[ $a == *f* ]]; then _g_deny "forced git push"; return 0; fi ;;
        esac
      done
      _g_ask "git push sends code to the remote"
      ;;
    reset)
      for a in "${r[@]}"; do
        [[ $a == --hard ]] && _g_deny "git reset --hard"
      done
      ;;
    clean)
      for a in "${r[@]}"; do
        if [[ $a == --force || ( $a == -[!-]* && $a == *f* ) ]]; then _g_deny "forced git clean"; fi
      done
      ;;
    filter-repo | filter-branch) _g_deny "git history rewrite" ;;
    apply | am) _g_deny "git $gs writes files: use the Edit or Write tool" ;;
    grep) _g_deny "git grep reads files: use the Grep tool" ;;
  esac
}

_g_rules_gh() {
  local a has_secret=0 has_repo=0 has_delete=0 has_auth=0 has_token=0
  local -a pos=()
  for a in "$@"; do
    case $a in
      secret | ssh-key) has_secret=1 ;;
      repo) has_repo=1 ;;
      delete) has_delete=1 ;;
      auth) has_auth=1 ;;
      token) has_token=1 ;;
    esac
    [[ $a == -* ]] || pos+=("$a")
  done
  if ((has_secret)) || ((has_repo && has_delete)) || ((has_auth && has_token)); then
    _g_deny "destructive gh command or one that exposes tokens"
    return 0
  fi
  case "${pos[0]-} ${pos[1]-}" in
    "run list" | "run view" | "run watch" | "pr checks") ;;
    *) _g_ask "gh command that is not a CI query" ;;
  esac
}

# _g_rules COMMAND PREVIOUS DEPTH ARGS... — per-command rules
_g_rules() {
  local cn=$1 prev=$2 depth=$3
  shift 3
  local -a args=("$@")
  local a sub="" j found=0 script="" collecting=0
  local -a subw=()
  case $cn in
    sed | tee | cp | mv | touch | install | ln | truncate | patch | rsync | tar | zip | unzip)
      _g_deny "$cn reads or writes files through the shell: use the Edit or Write tool (Read to view)"
      ;;
    cat | tac | head | tail | less | more | bat | nl | od | xxd | hexdump | strings | grep | egrep | \
      fgrep | rg | ag | find | fd | ls | tree | awk | gawk | mawk | cut | jq | sort | uniq | wc | \
      diff | cmp | stat | file)
      _g_deny "$cn reads files through the shell: use the Read, Grep or Glob tool"
      ;;
    curl | wget)
      for a in "${args[@]}"; do
        if [[ $a == --output* || $a == --remote-name* || $a == --output-document* ||
          $a == --output-file* ||
          ( $a == -[!-]* && $a == *[oOaA]* && $a != *O- ) ]]; then
          _g_deny "$cn writes a file: use the Write tool"
          break
        fi
      done
      ;;
  esac
  case $cn in
    sudo | su | doas | pkexec | chown | dd | shred | wipefs | mkfs | mkfs.*)
      _g_deny "forbidden command: $cn"
      ;;
    psql | pg_dump | pg_restore | pgcli | mysql | mysqldump | mariadb | mongosh | mongo | \
      mongodump | mongorestore | redis-cli | sqlcmd)
      _g_deny "database client: $cn"
      ;;
    rm)
      for a in "${args[@]}"; do
        if [[ $a == --recursive || $a == --force || $a == --no-preserve-root ||
          ( $a == -[!-]* && $a == *[rRf]* ) ]]; then
          _g_deny "recursive or forced rm"
        fi
      done
      ;;
    git) _g_rules_git "${args[@]}" ;;
    gh) _g_rules_gh "${args[@]}" ;;
    go)
      case ${args[0]-} in get | install) _g_ask "go ${args[0]} changes dependencies or installs binaries" ;; esac
      ;;
    npm | pnpm | yarn | bun)
      for a in "${args[@]}"; do
        [[ $a == -* ]] && continue
        sub=$a
        break
      done
      case $sub in
        publish | unpublish) _g_deny "package publishing ($cn)" ;;
        install | i | add | ci) _g_ask "$cn $sub installs dependencies" ;;
      esac
      ;;
    cargo) for a in "${args[@]}"; do [[ $a == publish ]] && _g_deny "package publishing (cargo)"; done ;;
    twine) for a in "${args[@]}"; do [[ $a == upload ]] && _g_deny "package publishing (twine)"; done ;;
    terraform | tofu)
      for a in "${args[@]}"; do
        [[ $a == apply || $a == destroy ]] && _g_deny "infrastructure change ($cn)"
      done
      ;;
    kubectl)
      for a in "${args[@]}"; do
        [[ $a == delete || $a == drain ]] && _g_deny "destructive kubectl"
      done
      ;;
    sh | bash | zsh | dash | ksh | ash | fish)
      for ((j = 0; j < ${#args[@]}; j++)); do
        a=${args[j]}
        if [[ $a == -[!-]* && $a == *c* ]]; then
          found=1
          script=${args[j + 1]-}
          break
        fi
        [[ $a == -* ]] || break
      done
      if ((found)); then
        _g_ask "$cn -c runs command text"
        _g_analyze_string "$script" $((depth + 1))
      elif [[ $prev == '|' ]]; then
        _g_deny "script read from a pipe ($cn)"
      elif ((has_hd)); then
        _g_analyze_string "$_HD_ALL" $((depth + 1))
      fi
      ;;
    python | python[0-9]* | pypy* | perl | ruby | irb | node | nodejs | deno | php | lua | tclsh | Rscript | \
      pwsh | powershell | osascript)
      _g_deny "$cn runs scripts that bypass the file tools: not allowed"
      ;;
    eval) _g_ask "eval runs command text" ;;
    find)
      for a in "${args[@]}"; do
        if ((collecting)); then
          if [[ $a == ';' || $a == '+' ]]; then
            collecting=0
            _g_subcommand "${subw[@]}"
            subw=()
          else
            subw+=("$a")
          fi
        else
          case $a in
            -delete) _g_ask "find -delete removes files" ;;
            -exec | -execdir | -ok | -okdir) collecting=1 ;;
          esac
        fi
      done
      ((collecting)) && ((${#subw[@]})) && _g_subcommand "${subw[@]}"
      ;;
  esac
}

# _g_subcommand WORDS... — analyzes an embedded command (find -exec)
_g_subcommand() {
  local s
  printf -v s '%s\x1f' "$@"
  _g_segment "$s" "" "$((${depth:-0} + 1))"
}

# _g_segment SEGMENT PREVIOUS DEPTH — one simple command
_g_segment() {
  local seg=$1 prev=$2 depth=$3
  local -a w=() a=() rt=() rtm=()
  local k x y op ci n cn mut=0 has_hd=0
  {
    local IFS=$_G_US
    # shellcheck disable=SC2206 # intentional split on US, with noglob on
    w=($seg)
    IFS=$' \t\n'
  }
  for ((k = 0; k < ${#w[@]}; k++)); do
    x=${w[k]}
    case ${x:0:1} in
      $'\x02') ;;
      $'\x01')
        op=${x:1}
        k=$((k + 1))
        if ((k < ${#w[@]})); then
          rt+=("${w[k]}")
          case $op in
            '<<' | '<<-' | '<<<') rtm+=(-1) ;;
            '<') rtm+=(0) ;;
            *) rtm+=(1) ;;
          esac
        fi
        [[ $op == '<<' || $op == '<<-' ]] && has_hd=1
        ;;
      *) a+=("$x") ;;
    esac
  done
  n=${#a[@]}
  ci=0
  while ((ci < n)); do
    x=${a[ci]}
    if [[ $x =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then
      ci=$((ci + 1))
      continue
    fi
    case $x in
      do | then | else | elif | if | while | until | '!' | '{' | time)
        ci=$((ci + 1))
        continue
        ;;
      env | command | builtin | exec | nohup | nice | setsid | stdbuf | timeout | xargs | ionice)
        ci=$((ci + 1))
        while ((ci < n)); do
          y=${a[ci]}
          case $y in
            --) ci=$((ci + 1)); break ;;
            -n | -I | -P | -d | -E | -L | -s | -a | -o | -e | -i | -u | -k | -c) ci=$((ci + 2)) ;;
            -*) ci=$((ci + 1)) ;;
            [0-9]*) [[ $x == timeout ]] && ci=$((ci + 1)); break ;;
            *) break ;;
          esac
        done
        continue
        ;;
    esac
    break
  done
  cn=""
  if ((ci < n)); then
    cn=${a[ci]}
    cn=${cn#\\}
    cn=${cn##*/}
    case $cn in
      cat | ls | head | tail | less | more | grep | egrep | fgrep | rg | wc | diff | stat | file | find | \
        tree | jq | bat | sort | uniq | cut | awk | realpath | basename | dirname | readlink | test | \
        '[' | '[[' | echo | printf | git | go | gofmt | cd | pwd | which | type | true | false) mut=0 ;;
      sed)
        for y in "${a[@]:ci+1}"; do
          [[ $y == --in-place* || ( $y == -[!-]* && $y == *i* ) ]] && mut=1
        done
        ;;
      *) mut=1 ;;
    esac
  fi
  for ((k = 0; k < n; k++)); do
    if ((k == ci)); then
      _g_scan_cmdword "${a[k]}"
    elif ((k < ci)); then
      _g_scan_word "${a[k]}" 0
    else
      _g_scan_word "${a[k]}" "$mut"
    fi
  done
  for ((k = 0; k < ${#rt[@]}; k++)); do
    ((rtm[k] >= 0)) || continue
    _g_scan_word "${rt[k]}" "${rtm[k]}"
    case ${rt[k]} in
      /dev/null | /dev/stdin | /dev/stdout | /dev/stderr) ;;
      *)
        if ((rtm[k] == 1)); then
          _g_deny "redirect writes a file: use the Write tool"
        else
          _g_deny "redirect reads a file: use the Read tool"
        fi
        ;;
    esac
  done
  if [[ -n $cn ]]; then
    _g_rules "$cn" "$prev" "$depth" "${a[@]:ci+1}"
    case $cn in curl | wget) _PIPE_FETCH=1 ;; esac
  fi
}

# _g_analyze_string TEXT DEPTH — analyzes a shell command text
_g_analyze_string() {
  local text=$1 depth=$2 i prev_sep=""
  local -a _SEGS=() _SEPS=()
  local _HD_ALL="" _PIPE_FETCH=0
  if ((depth > 4)); then
    _g_ask "commands nested too deeply"
    return 0
  fi
  _g_lex "$text"
  for i in "${!_SEGS[@]}"; do
    [[ $prev_sep == '|' ]] || _PIPE_FETCH=0
    _g_segment "${_SEGS[i]}" "$prev_sep" "$depth"
    prev_sep=${_SEPS[i]}
  done
  if ((${#_SEGS[@]} > 1)) || [[ ${_SEPS[0]-} == '&' ]]; then
    _g_deny "chained commands are not allowed: run one command per call"
  fi
}

# ------------------------------------------------------------------- API

_guard_decide_inner() {
  local tool=${1-} proj=${2-} p mut=0
  shift 2 2>/dev/null
  _G_PROJ=${proj%/}
  _G_CWD=${GUARD_CWD:-$_G_PROJ}
  _R_DENY=""
  _R_ASK=""
  _C_VAL=()
  _C_MUT=()
  case $tool in
    Bash)
      [[ -n ${GUARD_CWD-} ]] && _g_scan_value "$GUARD_CWD" 0 1
      [[ -n ${1-} ]] && _g_analyze_string "$1" 0
      ;;
    Read | Glob | Grep | Edit | Write | NotebookEdit)
      case $tool in Edit | Write | NotebookEdit) mut=1 ;; esac
      for p in "$@"; do
        _g_scan_value "$p" "$mut" 1
      done
      ;;
  esac
  _g_resolve
  if [[ -n $_R_DENY ]]; then
    GUARD_DECISION=deny
    GUARD_REASON=$_R_DENY
  elif [[ -n $_R_ASK ]]; then
    GUARD_DECISION=ask
    GUARD_REASON=$_R_ASK
  else
    GUARD_DECISION=allow
    GUARD_REASON=""
  fi
}

guard_decide() {
  local _noglob=0
  [[ $- == *f* ]] && _noglob=1
  set -f
  _guard_decide_inner "$@"
  ((_noglob)) || set +f
  return 0
}

# guard_log TOOL DECISION REASON PROJECT — one line in PROJECT/.claude/logs/guard.log.
# The reason never contains a secret's path or the command text.
guard_log() {
  local dir="$4/.claude/logs" ts reason=${3//[$'\n\r']/ }
  printf -v ts '%(%Y-%m-%dT%H:%M:%S%z)T' -1
  { mkdir -p "$dir" && printf '%s %s %s %s\n' "$ts" "$1" "$2" "$reason" >>"$dir/guard.log"; } 2>/dev/null || true
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
  guard_decide "$@"
  printf '%s\t%s\n' "$GUARD_DECISION" "$GUARD_REASON"
fi

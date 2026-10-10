#!/usr/bin/env bash
# Guard adapter for Antigravity (PreToolUse hook in .agents/hooks.json).
# Reads toolCall.name and toolCall.args from stdin, calls the core (scripts/guard-core.sh) and answers:
#   allow → nothing (Antigravity continues with its own permission review)
#   ask   → {"decision":"force_ask"}  (always asks, ignoring permissions already granted)
#   deny  → {"decision":"deny"}
# If anything fails, it asks the user.
# GUARD_PROJECT_DIR overrides the project directory (used by the test suite).

HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)

ask_now() {
  printf '{"decision":"force_ask","reason":"guard: %s"}\n' "$1"
  exit 0
}

command -v jq >/dev/null 2>&1 || ask_now "jq is missing, the call could not be analyzed"

fields=()
mapfile -d '' -t fields < <(jq --raw-output0 '
  (.toolCall.name // ""),
  (.workspacePaths[0] // ""),
  (.toolCall.args // {} | (
    (.Cwd // ""),
    (.CommandLine // .command // ""),
    ([to_entries[] | select((.key | test("path|file|dir"; "i")) and (.value | type == "string")) | .value] | join("\u001f"))
  ))' 2>/dev/null)
((${#fields[@]} >= 3)) || ask_now "hook input is unreadable"

name=${fields[0]}
workspace=${fields[1]}
cwd=${fields[2]}
command=${fields[3]-}
paths=${fields[4]-}

PROJECT=${GUARD_PROJECT_DIR:-${workspace:-$(cd "$HERE/../.." && pwd -P)}}
CORE="$PROJECT/scripts/guard-core.sh"
[[ -f $CORE ]] || CORE="$HERE/../../scripts/guard-core.sh"
[[ -f $CORE ]] || ask_now "scripts/guard-core.sh is missing"

case $name in
  run_command) tool=Bash ;;
  view_file | list_dir | find_by_name) tool=Read ;;
  grep_search) tool=Grep ;;
  write_to_file) tool=Write ;;
  replace_file_content | multi_replace_file_content) tool=Edit ;;
  *) exit 0 ;;
esac

inputs=()
if [[ $tool == Bash ]]; then
  inputs=("$command")
else
  IFS=$'\x1f' read -r -a inputs <<<"$paths"
fi

# shellcheck source-path=SCRIPTDIR source=../../scripts/guard-core.sh
source "$CORE" || ask_now "failed to load the core"
[[ -n $cwd ]] && export GUARD_CWD=$cwd
guard_decide "$tool" "$PROJECT" "${inputs[@]}"

case $GUARD_DECISION in
  allow) exit 0 ;;
  deny)
    guard_log "$tool" deny "$GUARD_REASON" "$PROJECT"
    jq -cn --arg r "guard: action denied: $GUARD_REASON" '{decision:"deny", reason:$r}' ||
      printf '{"decision":"deny"}\n'
    ;;
  *)
    guard_log "$tool" ask "$GUARD_REASON" "$PROJECT"
    jq -cn --arg r "guard: $GUARD_REASON" '{decision:"force_ask", reason:$r}' ||
      ask_now "confirmation required"
    ;;
esac

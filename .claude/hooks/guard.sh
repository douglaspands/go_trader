#!/usr/bin/env bash
# Guard adapter for Claude Code (PreToolUse hook).
# Reads the hook JSON from stdin, calls the core (scripts/guard-core.sh) and answers:
#   allow → nothing (the decision is left to the rules in settings.json)
#   ask   → hookSpecificOutput.permissionDecision = "ask"
#   deny  → exit 2 with the reason on stderr
# If anything fails (jq missing, invalid JSON, core missing), it asks the user.
# GUARD_PROJECT_DIR overrides the project directory (used by the test suite).

HERE=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)
PROJECT=${GUARD_PROJECT_DIR:-${CLAUDE_PROJECT_DIR:-$(cd "$HERE/../.." && pwd -P)}}
CORE="$PROJECT/scripts/guard-core.sh"
[[ -f $CORE ]] || CORE="$HERE/../../scripts/guard-core.sh"

ask_now() {
  printf '{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"ask","permissionDecisionReason":"guard: %s"}}\n' "$1"
  exit 0
}

command -v jq >/dev/null 2>&1 || ask_now "jq is missing, the call could not be analyzed"
[[ -f $CORE ]] || ask_now "scripts/guard-core.sh is missing"

fields=()
mapfile -d '' -t fields < <(jq --raw-output0 '
  .tool_name,
  (.cwd // ""),
  (.tool_input // {} | (
    (.command // empty),
    (.file_path // empty),
    (.path // empty),
    (.notebook_path // empty),
    ((.pattern // empty) | select(type == "string" and test("^(/|~|\\.\\./)")))
  ))' 2>/dev/null)
((${#fields[@]} >= 2)) || ask_now "hook input is unreadable"

tool=${fields[0]}
cwd=${fields[1]}
inputs=("${fields[@]:2}")

case $tool in
  Bash | Read | Glob | Grep | Edit | Write | NotebookEdit) ;;
  *) exit 0 ;;
esac

# shellcheck source-path=SCRIPTDIR source=../../scripts/guard-core.sh
source "$CORE" || ask_now "failed to load the core"
[[ -n $cwd ]] && export GUARD_CWD=$cwd
guard_decide "$tool" "$PROJECT" "${inputs[@]}"

case $GUARD_DECISION in
  allow) exit 0 ;;
  deny)
    guard_log "$tool" deny "$GUARD_REASON" "$PROJECT"
    echo "guard: action denied: $GUARD_REASON" >&2
    exit 2
    ;;
  *)
    guard_log "$tool" ask "$GUARD_REASON" "$PROJECT"
    jq -cn --arg r "guard: $GUARD_REASON" \
      '{hookSpecificOutput:{hookEventName:"PreToolUse", permissionDecision:"ask", permissionDecisionReason:$r}}' ||
      ask_now "confirmation required"
    ;;
esac

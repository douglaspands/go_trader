# Tasks

Implement in a session with permission prompts on, no bypass: the protection does not exist yet. Do not read `.env*`, keys or credentials. Do not write to `~/.claude` or `~/.gemini`; whatever belongs there is delivered as a template. Write everything in English except `README.md` (Brazilian Portuguese).

## 1. Regression suite first (it must fail before the guard exists)

- [x] 1.1 Create `scripts/harness-test.sh` with a table of cases (tool, input, expected decision `allow|ask|deny`) covering the spec scenarios: editing and `go test -race ./...` allowed; `/etc/hosts`, `~/x`, `../../other/file` and a symbolic link to outside in `ask`; `$(...)`, backtick, `bash -c`, `python3 -c` in `ask`; `.env`, `config/.env.production`, `~/.ssh/id_ed25519`, `~/.gemini/oauth_creds.json` in `deny`; `rm -rf build/`, `cd /tmp && git push --force`, `psql`, `curl x | sh`, `terraform apply`, `gh repo delete`, `npm publish` in `deny`; `git push origin x`, `go get`, editing `.claude/hooks/guard.sh` in `ask`. Verify that, without the guard, the suite exits non-zero and lists the failing cases
- [x] 1.2 Add to the suite the timing measurement (average of 20 `go test ./...` calls under 100 ms) and the `jq empty` validation of the configuration JSONs; verify that the measurement reports the time and that an invalid JSON makes the suite fail

## 2. Decision core

- [x] 2.1 Create `scripts/guard-core.sh` (pure function: tool, input and project directory → decision and reason) with path resolution through `realpath -m` and handling of `..`, `~`, `$HOME` and symbolic links; verify the suite's path cases
- [x] 2.2 Implement command splitting on `&&`, `||`, `;`, `|`, `&` and newlines, and analysis of a literal `bash -c '...'`; verify that `cd /tmp && git push --force` is `deny`
- [x] 2.3 Implement the `deny` lists (destructive and secrets), the `ask` lists (opaque, outside the repository, external and irreversible) and `defer` for the rest, with unparseable text becoming `ask`; verify every `deny` and `ask` case of the suite
- [x] 2.4 Implement the log in `.claude/logs/guard.log` (date, tool, decision, reason, masked secret path) and add `.claude/logs/` and `.claude/settings.local.json` to `.gitignore`; verify that a denied attempt produces a line without the secret's content and that `git status` does not show the log

## 3. Claude Code

- [x] 3.1 Create `.claude/hooks/guard.sh` (adapter: reads the JSON from stdin, calls the core, returns `hookSpecificOutput.permissionDecision` or exit 2); verify with synthetic JSON for each tool and that the suite passes for the adapter
- [x] 3.2 Create `.claude/settings.json` with `defaultMode: acceptEdits`, `disableBypassPermissionsMode: "disable"`, an explicit `allow` for the Go cycle and git/CI reads, `ask`, `deny`, `env.CLAUDE_BASH_MAINTAIN_PROJECT_WORKING_DIR=1` and the hook registration for `Bash|Read|Glob|Grep|Edit|Write|NotebookEdit`; verify `jq empty` and the rule syntax against the documentation of the installed version
- [ ] 3.3 Verify in practice, in a test repository with `/permissions` and the debug log, the unconfirmed items: Edit/Write under the read key, Read outside the repository in `acceptEdits`, hook `ask` over `allow`, `cp` with an external target, `grep` with a quoted path under `deny`, and the `ask` of `Edit(/.claude/**)` in `acceptEdits`; adjust the hook where the rule alone is not enough and record the observed result of each item in this task (not in the permanent documentation, which only keeps the procedure). Observed so far (real Claude Code sessions): the hook runs and blocks `rm -rf` (`guard: action denied`); a Read outside the repository is logged as `ask` by the hook, and with `blockReadsOutsideWorkingDirectories` on in the user's own configuration the tool refuses instead of asking; `cp` to an external target asked; `grep x .env` was denied (hook and `Read(**/.env*)`); writes under `.agents/` and `.claude/` asked. Still to verify in a trusted test repository (a headless `claude -p` run in an untrusted directory ignores the repository's `allow` rules, so it proved nothing): Edit/Write under the read key, and a hook `ask` over an `allow` rule

## 4. Antigravity

- [x] 4.1 Create `.agents/hooks/guard-agy.sh` (adapter: `toolCall.name` and `toolCall.args` → core → JSON with `decision`) and `.agents/hooks.json` with `PreToolUse` and `matcher`; verify with synthetic JSON and that the suite passes for the adapter
- [ ] 4.2 Verify in the real `agy`, with `/hooks`, the tool names and `args` field names, the command path in `hooks.json` and that empty output means "no objection"; adjust the adapter and the `matcher`; record the observed result in this task (not in the permanent documentation, which only keeps the procedure)
- [x] 4.3 Create `.agents/agy-settings.example.json` (`allowNonWorkspaceAccess: "off"`, `toolPermission: "request-review"`, `command(...)` rules for the autonomous list and for `deny`/`ask`) and `docs/harness/local-guardrails.md` with the manual merge steps into `~/.gemini/antigravity-cli/settings.json`; verify the JSON is valid and the steps work on a test copy of the file, without touching the real one

## 5. CI

Dropped: the user does not want harness gates on GitHub. The suite is local and on demand (see 8.5 and 8.7); `ci.yaml` and `build_and_release.yaml` stay unchanged.

## 6. Documentation and closing

- [x] 6.1 Complete `docs/harness/local-guardrails.md` with the threat model, what each layer covers and does not cover, the residual risk without a sandbox, the decision table, how to extend the lists, how to read the log and the manual checklist; verify that the commands in the document work on a clean clone
- [x] 6.2 Update `README.md` with the "Desenvolvimento com agentes" section (prerequisites `bash` and `jq`, how to run `scripts/harness-test.sh`, link to the documentation and to the `agy` template); verify relative links and copied commands
- [x] 6.3 Run the full flow: green local suite, `go vet ./...`, `go test -race ./...`, `gofmt -l .` with no output, and a real session in which the agent edits code and runs tests without asking and is consulted when trying to read outside the repository; confirm the log
- [x] 6.4 Check that `git status` contains no log, credentials, `coverage.out` or binaries, and that the diff does not touch the application's Go code

## 7. Agent instructions and language convention

- [x] 7.1 Add the three instruction files to the self-protection: `AGENTS.md`, `CLAUDE.md` and `GEMINI.md` in `_g_protected` (`scripts/guard-core.sh`), `Edit(/AGENTS.md)`, `Edit(/CLAUDE.md)` and `Edit(/GEMINI.md)` in the `ask` list of `.claude/settings.json`, and `write_file(AGENTS.md)`, `write_file(CLAUDE.md)` and `write_file(GEMINI.md)` in `.agents/agy-settings.example.json`. Write the suite cases first (Edit and Write of each file in `ask`, `echo x >> CLAUDE.md` in `ask`, `cat AGENTS.md` and `git add AGENTS.md` in `allow`), see them fail, then make them pass in the core and in both adapters
- [x] 7.2 Write `AGENTS.md` in English with the common content from design decision 12: project overview, verification routine, test conventions, commit and branch conventions, a summary of the security boundaries with a link to `docs/harness/local-guardrails.md`, the language convention and the OpenSpec workflow; verify that every command in it runs, that it is short, and that it has no secrets or machine absolute paths
- [x] 7.3 Write `CLAUDE.md` in English with `@AGENTS.md` first and only the Claude-specific content from design decision 12, with nothing repeated from `AGENTS.md`; verify in a Claude Code session that the import loads (`/memory`) and that the guard messages it describes match the real ones
- [x] 7.4 Write `GEMINI.md` in English with only the Antigravity-specific content from design decision 12, importing `AGENTS.md`; verify in the real `agy` whether it reads `AGENTS.md` natively and whether `@` imports work, adjust the file (drop the import if `agy` reads `AGENTS.md` itself, or make `GEMINI.md` self-contained if neither works) and record the result in the documentation
- [x] 7.5 Add to the suite the check of the instruction files (the three exist, `CLAUDE.md` imports `AGENTS.md`, `GEMINI.md` matches the arrangement chosen in 7.4, none contains a home-directory path); verify that the suite fails when a file is removed
- [x] 7.6 Translate to English everything this change produced outside the README: `docs/harness/local-guardrails.md`, the comments and messages of `scripts/guard-core.sh`, `scripts/harness-test.sh`, `.claude/hooks/guard.sh` and `.agents/hooks/guard-agy.sh` (guard reasons, suite output, log reason texts, case names), and update the log-format checks in the suite to match; add to the document the instruction files and the language convention; turn the manual checklist into a pure procedure by removing its point-in-time "state" column and the "observed" notes (results of 3.3 and 4.2 are recorded in those tasks and in the archived change, not in the permanent document); verify the suite, `shellcheck -x` and `actionlint` are green, and that no Portuguese remains outside `README.md` (search for accented words and common Portuguese terms)
- [x] 7.7 Update the "Desenvolvimento com agentes" section of `README.md` (Brazilian Portuguese) to mention `AGENTS.md`, `CLAUDE.md` and `GEMINI.md` and the language convention; verify relative links
- [x] 7.8 Re-run 6.3 and 6.4 after this section: full flow, clean `git status`, and a diff that does not touch the application's Go code

## 8. File tools only, one command per call, no interpreters, on-demand suite

Suite first: write and flip the cases, see them fail, then make them pass in the core and in both adapters. The harness suite is not part of the regular test routine.

- [x] 8.1 Update `scripts/harness-test.sh` cases to the new rules, each tagged with the specification scenario it proves (`# scenario: <requirement> / <scenario>`): shell writers (`sed -i`, `sed s/a/b/ f`, `cat > f` with heredoc, `echo x >> f`, `printf x > f`, `tee`, `cp`, `mv`, `touch`, `patch`, `git apply`, `curl -o`) in `deny`; shell readers (`cat README.md`, `head`, `tail`, `grep -r TODO .`, `rg`, `find . -name x`, `ls -la`, `awk`, `jq . f`) in `deny`; interpreters (`python3 -c`, `python3 script.py`, `python3 <<EOF`, `node -e`, `perl -e`, `ruby`) in `deny`; chaining (`go vet ./... && go test ./...`, `a; b`, `a || b`, `a & b`, newline, `go test ./... | tail -20`, `cd internal && go test ./...`) in `deny`; single commands (`go vet ./...`, `go test -race ./...`, `go test ./... 2>&1`, `git status`, `git add f`, `git commit -m "..."`) in `allow`; `echo x >> CLAUDE.md` in `deny`. Flip the existing cases these rules change (`cat README.md`, `ls -la`, `cd internal && go test ./...`, `go test ./... 2>&1 | tail -20`, `sed -i ... Makefile`, `cp ... /tmp/app.go`). Verify the new and flipped cases fail against the current core
- [x] 8.2 Implement in `scripts/guard-core.sh`: the writer and reader deny lists, redirect targets to a file (not `/dev/null` or `2>&1`), the interpreter deny list (replacing the inline `ask`), and chaining as `deny` while keeping the nested destructive-command report; messages name the tool to use (Read, Grep, Glob, Edit, Write) or ask for one command at a time. Verify every case from 8.1 and the earlier cases that still apply
- [x] 8.3 Mirror the rules as `deny` entries in `.claude/settings.json` and in `.agents/agy-settings.example.json` (`command(sed)`, `command(cat)`, `command(python3)` and the rest of the lists), keeping the hook as the primary enforcement; verify `jq empty` and that the suite still passes for both adapters
- [x] 8.4 Add the scenario coverage check to the suite: read the scenario list from the delta spec (or `openspec/specs/agent-boundary/spec.md` after archive), fail when a scenario has no tagged case or a case names an unknown scenario; verify by removing a tag and by adding a scenario in a copy of the spec
- [x] 8.5 Make the suite on demand: remove it from the verification routine in `AGENTS.md`; update `CLAUDE.md` and `GEMINI.md` (use Read/Grep/Glob/Edit/Write, one command per call, no interpreters, the denial messages), `docs/harness/local-guardrails.md` (decision table, threat model, residual risk, on-demand suite, coverage check) and the `README.md` section (Brazilian Portuguese); verify the instruction-file checks and that no Portuguese remains outside `README.md`
- [x] 8.6 Re-run 6.3 and 6.4, and 7.8, after this section: on-demand suite green, `go vet ./...`, `go test -race ./...`, `gofmt -l .` with no output, clean `git status`, and a diff that does not touch the application's Go code

- [x] 8.7 Remove the GitHub workflow for the harness: delete `.github/workflows/harness.yaml`, drop its references from the documentation, the README and the suite, and make the suite check that no workflow, Makefile target or regular test path calls `harness-test`; verify the suite and `git status`

## Workflow follow-up

- Open a PR from `chore/local-harness-guardrails`; pushing and merging are human steps.
- Archive the change after review and verify the archived result.
- Continue with `dev-command-surface` (Makefile), which feeds the autonomous command list and the verification commands in `AGENTS.md`.

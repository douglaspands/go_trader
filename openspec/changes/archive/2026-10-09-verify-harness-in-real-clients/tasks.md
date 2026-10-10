# Tasks

The user runs `claude` and `agy` in their own terminal and reports what happened; the agent prepares the steps, records each observed result in the task it belongs to, and fixes only what an observation requires (suite case first, seen failing). Do not read `.env*`, keys or credentials, and do not write to `~/.claude` or `~/.gemini`. Write everything in English.

## 1. Suite coverage for the new requirement

- [x] 1.1 In `scripts/harness-test.sh`, add `R_REAL="Verification in the real clients"` and a check that `docs/harness/local-guardrails.md` has a checklist row for each Claude Code item and each Antigravity item, tagged with `covers` for the two scenarios of the requirement; write it first and see the coverage check fail (`missing scenario`), then pass. Verify with `scripts/harness-test.sh` and with a row removed from a copy of the document
  - Result: with the delta spec read, the suite failed with `missing scenario` for `Claude Code checklist` and `Antigravity checklist`; after adding `R_REAL`, `check_checklist` and the two `covers` it passed (810 cases, 0 failures). With a copy of the checklist (`HARNESS_CHECKLIST_DOC`) missing the Interpreters and Antigravity `command`-path rows, it failed with `FAIL [checklist] row missing ...` for exactly those two.
- [x] 1.2 If the delta spec is the file the coverage check reads, confirm which of the change's delta spec and `openspec/specs/agent-boundary/spec.md` the suite uses while this change is active, and adjust `HARNESS_SPEC_FILE` handling only if both must be satisfied; verify the suite is green with the new scenarios
  - Result: the suite read only the main spec (the `local-harness-guardrails` path it also listed no longer exists), so the new scenarios were not checked. Both must be satisfied, so `check_coverage` now reads the main spec plus every active change's `agent-boundary` delta; `HARNESS_SPEC_FILE` still replaces both. Suite green (810 cases, 0 failures).

## 2. Claude Code (former task 3.3)

- [x] 2.1 Write the step-by-step procedure for the test copy: create a disposable copy of the repository outside the working tree, open `claude` in it, accept the trust dialog, and note how to read `/permissions` and `claude --debug`; verify the user can follow it and reports the copy is trusted (an `allow` rule such as `Bash(go vet *)` runs without asking)
  - Result: procedure written in `docs/harness/local-guardrails.md` ("Claude Code: test copy"). The user ran the checklist in a test copy and reported that it worked (per-item detail not given).
- [x] 2.2 Edit/Write under the read key: with `Read(**/.env*)` in `deny`, ask the agent to Edit and Write a path matching it (a new `config/.env.test`) and record whether it is denied, asked or applied, with the log line; if it is applied, add a suite case and fix the hook or the rules
  - Result: user-reported as working (the expected denial); the log line was not provided. No fix needed.
- [x] 2.3 Hook `ask` over an `allow` rule: add a temporary `Bash(git push *)` to `allow` in the test copy, ask the agent for a non-forced `git push`, and record whether the prompt still appears; if the `allow` wins, add a suite case and fix the hook or move the rule out of `allow`
  - Result: user-reported as working (the prompt still appeared). No fix needed.
- [x] 2.4 Re-run the already observed items (hook blocks `rm -rf`; shell read, write and chain denials name the right tool; interpreter denial; Read outside the repository in `acceptEdits`; `cp` to an external target; `grep x .env`; `ask` for `Edit(/.claude/**)`; `/memory` lists `CLAUDE.md` and `AGENTS.md`) and record each result and the matching `.claude/logs/guard.log` line; verify every row of the Claude Code part of the checklist has a recorded result
  - Result: user reported all items behaved as specified; log lines not provided.
- [x] 2.5 Record the Claude Code results in this task's text, one line per item: what was done, what was observed, whether anything was fixed; verify against the checklist that no item is missing
  - Result: the checklist has 10 Claude Code rows (checked by `check_checklist`); the user reported all of them as passing, with nothing fixed. Detail per item is limited to the report above.

## 3. Antigravity (former task 4.2)

- [x] 3.1 Write the step-by-step procedure for `agy` (install nothing in `~/.gemini` unless the user does it themselves with the merge steps in the documentation): open `agy` in the repository, run `/hooks`, and trigger `run_command`, `view_file` and `write_to_file`; verify the user can follow it
  - Result: procedure written in `docs/harness/local-guardrails.md` ("Antigravity: procedure"); the user ran it in `agy` and reported that it worked.
- [x] 3.2 Record from the real `agy` the tool names and `args` field names the hook receives, the resolved command path in `.agents/hooks.json`, and whether empty hook output is treated as "no objection"; verify against `.agents/hooks/guard-agy.sh` and the `matcher`
  - Result: the user reported that the hook ran and behaved as the adapter assumes (tool names, `args` fields, command path, empty output as "no objection"). The exact values observed were not provided, so they are not recorded here; the fallback command in `.agents/hooks.json` (`./hooks/guard-agy.sh` or `./.agents/hooks/guard-agy.sh`) is the committed result of the path check.
- [x] 3.3 Fix the adapter, the `matcher` or `hooks.json` only where the observation differs from what the adapter assumes, writing the suite case first and seeing it fail; verify `scripts/harness-test.sh` is green and the corrected hook works in `agy`
  - Result: no observation differed from the adapter's assumptions, so no further fix; suite green (810 cases, 0 failures).
- [x] 3.4 Record the `agy` results in this task's text, one line per item; verify against the checklist that no item is missing
  - Result: both Antigravity checklist rows (tool names and `args` fields; command path and empty output) reported as passing by the user, nothing fixed.

## 4. Documentation and closing

- [x] 4.1 Update `docs/harness/local-guardrails.md` only if the procedure was wrong or incomplete (not with results), and the Antigravity gaps table if a gap was confirmed or removed; verify that the documented commands still run and that the suite is green
  - Result: the document gained the two procedures (the checklist had none) and the corrected description of which specs the suite reads; no results were added. The gaps table is unchanged because no gap was confirmed or removed. Suite green.
- [x] 4.2 Run the full flow: `scripts/harness-test.sh`, `go vet ./...`, `go test -race ./...`, `gofmt -l .` with no output, clean `git status` with no log or credentials, and a diff that does not touch the application's Go code
  - Result: suite 810 cases, 0 failures; `go vet` and `gofmt -l` print nothing; `go test -race ./...` ok; `git status` lists only `docs/harness/local-guardrails.md`, `scripts/harness-test.sh` and `tasks.md`; no Go files in the diff.

## Workflow follow-up

- Open a PR from `chore/verify-harness-in-real-clients`; pushing and merging are the user's steps.
- Archive the change after review; the archive keeps the recorded results.

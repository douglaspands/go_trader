# Tasks

The user runs `claude` and `agy` in their own terminal and reports what happened; the agent prepares the steps, records each observed result in the task it belongs to, and fixes only what an observation requires (suite case first, seen failing). Do not read `.env*`, keys or credentials, and do not write to `~/.claude` or `~/.gemini`. Write everything in English.

## 1. Suite coverage for the new requirement

- [ ] 1.1 In `scripts/harness-test.sh`, add `R_REAL="Verification in the real clients"` and a check that `docs/harness/local-guardrails.md` has a checklist row for each Claude Code item and each Antigravity item, tagged with `covers` for the two scenarios of the requirement; write it first and see the coverage check fail (`missing scenario`), then pass. Verify with `scripts/harness-test.sh` and with a row removed from a copy of the document
- [ ] 1.2 If the delta spec is the file the coverage check reads, confirm which of the change's delta spec and `openspec/specs/agent-boundary/spec.md` the suite uses while this change is active, and adjust `HARNESS_SPEC_FILE` handling only if both must be satisfied; verify the suite is green with the new scenarios

## 2. Claude Code (former task 3.3)

- [ ] 2.1 Write the step-by-step procedure for the test copy: create a disposable copy of the repository outside the working tree, open `claude` in it, accept the trust dialog, and note how to read `/permissions` and `claude --debug`; verify the user can follow it and reports the copy is trusted (an `allow` rule such as `Bash(go vet *)` runs without asking)
- [ ] 2.2 Edit/Write under the read key: with `Read(**/.env*)` in `deny`, ask the agent to Edit and Write a path matching it (a new `config/.env.test`) and record whether it is denied, asked or applied, with the log line; if it is applied, add a suite case and fix the hook or the rules
- [ ] 2.3 Hook `ask` over an `allow` rule: add a temporary `Bash(git push *)` to `allow` in the test copy, ask the agent for a non-forced `git push`, and record whether the prompt still appears; if the `allow` wins, add a suite case and fix the hook or move the rule out of `allow`
- [ ] 2.4 Re-run the already observed items (hook blocks `rm -rf`; shell read, write and chain denials name the right tool; interpreter denial; Read outside the repository in `acceptEdits`; `cp` to an external target; `grep x .env`; `ask` for `Edit(/.claude/**)`; `/memory` lists `CLAUDE.md` and `AGENTS.md`) and record each result and the matching `.claude/logs/guard.log` line; verify every row of the Claude Code part of the checklist has a recorded result
- [ ] 2.5 Record the Claude Code results in this task's text, one line per item: what was done, what was observed, whether anything was fixed; verify against the checklist that no item is missing

## 3. Antigravity (former task 4.2)

- [ ] 3.1 Write the step-by-step procedure for `agy` (install nothing in `~/.gemini` unless the user does it themselves with the merge steps in the documentation): open `agy` in the repository, run `/hooks`, and trigger `run_command`, `view_file` and `write_to_file`; verify the user can follow it
- [ ] 3.2 Record from the real `agy` the tool names and `args` field names the hook receives, the resolved command path in `.agents/hooks.json`, and whether empty hook output is treated as "no objection"; verify against `.agents/hooks/guard-agy.sh` and the `matcher`
- [ ] 3.3 Fix the adapter, the `matcher` or `hooks.json` only where the observation differs from what the adapter assumes, writing the suite case first and seeing it fail; verify `scripts/harness-test.sh` is green and the corrected hook works in `agy`
- [ ] 3.4 Record the `agy` results in this task's text, one line per item; verify against the checklist that no item is missing

## 4. Documentation and closing

- [ ] 4.1 Update `docs/harness/local-guardrails.md` only if the procedure was wrong or incomplete (not with results), and the Antigravity gaps table if a gap was confirmed or removed; verify that the documented commands still run and that the suite is green
- [ ] 4.2 Run the full flow: `scripts/harness-test.sh`, `go vet ./...`, `go test -race ./...`, `gofmt -l .` with no output, clean `git status` with no log or credentials, and a diff that does not touch the application's Go code

## Workflow follow-up

- Open a PR from `chore/verify-harness-in-real-clients`; pushing and merging are the user's steps.
- Archive the change after review; the archive keeps the recorded results.

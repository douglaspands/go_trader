# Design

## Context

See `proposal.md` for the motivation. The user's decision is to run agents on the host, with no container and no sandbox, and a single limit: the repository directory. Leaving it requires permission. Checked state and documentation:

- **Claude Code 2.1.296** (official documentation, with some items unconfirmed below). There is no single setting that gives "outside the repository asks, inside is autonomous":
  - In `default` mode, Read, Grep and Glob outside the working directory ask, and Edit/Write always ask. In `acceptEdits`, edits inside the project are automatic and edits outside ask.
  - Shell read commands (`ls`, `cat`, `grep`, `find`, `cd`, read-only `git`) run **without asking on any path**, including `cat /etc/passwd` and `cat ../x`, unless `permissions.blockReadsOutsideWorkingDirectories` is on.
  - That setting makes the file tools **refuse** reads outside the repository (it is not a prompt; the agent is told to ask for `/add-dir`) and makes those shell commands ask.
  - `deny` rules for `Read(...)` also cover `cat`, `head`, `tail`, `sed`, `tee` and redirects that name the file, but not scripts that open the file without naming it.
  - A `PreToolUse` hook can return `permissionDecision: "ask"`, which forces the prompt, and receives `cwd` on stdin. `auto` mode does not fit: whatever a rule does not resolve goes to a classifier, not to the user.
  - `CLAUDE.md` is the project instructions file; it does not read `AGENTS.md` by itself but supports `@path` imports.
- **Antigravity `agy` 1.3.2** (official documentation): `allowNonWorkspaceAccess` (default off) restricts reads and writes to the project; `command(prefix)`, `read_file`, `write_file` rules with `deny > ask > allow` and no path glob; `PreToolUse` hooks with `decision` `allow|ask|force_ask|deny|deny_unless_prior_grant` in `.agents/hooks.json` (project level). **The CLI's permission configuration only exists at the global level** (`~/.gemini/antigravity-cli/settings.json`), so it cannot be versioned. Whether `agy` reads `AGENTS.md` by itself, and whether `GEMINI.md` accepts imports, is not confirmed in the official documentation (third-party sources say `GEMINI.md` wins over `AGENTS.md` when both exist).
- The user's environment already has `blockReadsOutsideWorkingDirectories` on in their own configuration (reads outside the repository during diagnosis were refused). A project file can turn this key on, but cannot turn it off.
- The host has `bash`, `jq`, `rg`; `bwrap` works, but the user decided **not** to use a sandbox.
- The repository has no database and no production environment variables; the risk comes from the host and the user's credentials.

## Goals / Non-Goals

**Goals:**
- Full autonomy inside the repository for the Go cycle, without repetitive permission prompts.
- Any access outside the repository goes through the user; secrets and destructive commands are denied.
- A fast, tested guard with the same logic in both harnesses.
- Agents start every session with the project's shared rules, plus only what is specific to their harness.

**Non-Goals:**
- Kernel-level containment (sandbox, container, namespaces): out of scope by the user's decision. It can be turned on later (`sandbox.enabled`) without redoing this change.
- Makefile, lint, Node/npm and frontend (following changes). The Makefile change will feed the command list in `AGENTS.md`.
- Changing files outside the repository, including `~/.claude` and `~/.gemini`. Whatever must live there is delivered as a template for the user to install.

## Decisions

**1. Three layers of defense, no OS isolation.**
(a) Permission rules in `.claude/settings.json`; (b) a guard hook that inspects the call; (c) a regression suite. Without a sandbox, none of them is a kernel boundary: a program the agent runs (`go test`, `go run`, `make`) can read anything the user can. This is an **accepted residual risk** recorded in the documentation; the layers reduce the chance of the agent getting there unseen.

**2. The hook covers every tool, not only Bash.**
`matcher` is `Bash|Read|Glob|Grep|Edit|Write|NotebookEdit`. For file tools the guard resolves the path (with `realpath -m`, to catch `..` and symbolic links) and returns `ask` if it is outside the repository. This makes the behavior deterministic instead of depending on each mode's default (whether Read outside the repository asks in `acceptEdits` is not confirmed in the documentation). The repository is the hook's `cwd` (`CLAUDE_PROJECT_DIR` as reinforcement).

**3. For Bash, decide in order: deny, ask, defer.**
1. **Deny** (`deny`, exit 2 / `decision: deny`): destructive and secret patterns from the spec, matched on the real command words after splitting on `&&`, `||`, `;`, `|`, `&` and newlines, and inside `bash -c '...'` when the text is literal. Quoted text is data, not a command.
   The `deny` set also holds the rules of decisions 14, 15 and 16: shell file reads and writes, chaining and scripting interpreters.
2. **Ask** (`ask`): command substitution, backticks, `eval`, `sh -c`/`bash -c`, and any token that resolves outside the repository (absolute, `~`, `$HOME`, `$XDG_*`, `..` that escapes, symbolic link).
3. **Defer** (no output): the result falls to the `allow`/`ask`/`deny` rules in `settings.json`.
The hook never returns `allow`, so it does not override an `ask` or `deny` rule (confirmed in the documentation: a hook `allow` does not beat those rules, but keeping it this way avoids depending on that). Text the guard cannot parse becomes `ask`, never `allow`.
Rejected alternative: `Bash(...)` rules only. The documentation says they do not match `bash -c`, an absolute binary path or variables, and that they are not a security boundary.

**4. Implementation in bash + `jq`, no Python.**
Each call costs a few short processes, which fits in 100 ms. The decision logic lives in `scripts/guard-core.sh` (a pure function: tool, command or path and the project directory in; decision and reason out). Two thin adapters read each harness's JSON and write the response in its format: `.claude/hooks/guard.sh` (Claude: `hookSpecificOutput.permissionDecision`) and `.agents/hooks/guard-agy.sh` (Antigravity: `decision`, `toolCall.name` and `toolCall.args`; the tool names and `args` field names come from the documentation and are `[VERIFY]` on first real run). Parity comes from both calling the same core.

**5. Explicit autonomous command list, no broad wildcard.**
In `settings.json`, `allow` has `Bash(go build *)`, `go test`, `go vet`, `go fmt`, `go list`, `go mod tidy`, `go mod download`, `gofmt`, `git status|diff|log|show|branch|switch|add|commit|stash`, `gh run list|view|watch`, `gh pr checks` and, when they exist, explicit Makefile targets. No `Bash(*)`, `make *` or `npx *`, since a Makefile or package edited by the agent would become authorized arbitrary code. `defaultMode` is `acceptEdits`. Never `bypassPermissions` (`disableBypassPermissionsMode: "disable"`). `ask` for `git push`, `go get`, `go install`, `npm install` and for editing `.claude/`, `.agents/`, `.github/workflows/`, `Makefile`, `AGENTS.md`, `CLAUDE.md` and `GEMINI.md`. `deny` for the spec's patterns, reinforcing the hook.

**6. `CLAUDE_BASH_MAINTAIN_PROJECT_WORKING_DIR=1`.**
Documented as a way to start every command in the project directory, which gives the guard a predictable `cwd` and avoids a persistent `cd` outside the repository.

**7. The project does not turn on `blockReadsOutsideWorkingDirectories`.**
With the key on, the file tools refuse instead of asking, which contradicts the request ("ask for permission"). The project relies on the hook to ask. The user's global configuration still applies and, if they want asking instead of refusing, they must change it in their own file; this goes in the documentation.

**8. Antigravity: versioned hook + global template installed by the user.**
`.agents/hooks.json` with `PreToolUse` and a `matcher` over `run_command|view_file|list_dir|find_by_name|grep_search|write_to_file|replace_file_content|multi_replace_file_content`, calling `guard-agy.sh`. The adapter asks with `force_ask` (so a previously granted "always allow" does not skip the prompt) and denies with `deny`; on "no objection" it prints nothing. `.agents/agy-settings.example.json` carries `allowNonWorkspaceAccess: "off"`, `toolPermission: "request-review"` and `command(...)` rules equivalent to the autonomous list. The agent does **not** write to `~/.gemini`; the user merges the template manually. Antigravity gaps: no path glob (the hook compensates); global user rules and hooks may add to the project's.

**9. Self-protection of the configuration.**
Editing `.claude/**`, `.agents/**`, `.github/workflows/**`, `Makefile`, `scripts/guard*.sh`, `scripts/harness-test.sh`, `AGENTS.md`, `CLAUDE.md` and `GEMINI.md` is `ask`, even in `acceptEdits`, so the agent cannot weaken the guard or rewrite its own instructions without the user seeing it. The hook covers this itself for file tools and for shell writes (redirects, `sed -i`, `cp`, `rm`…), in case a rule's `ask` does not beat accept-edits mode (unconfirmed).

**10. Logging.**
`.claude/logs/guard.log` (format `date tool decision reason`), in `.gitignore`. The log holds the reason and kind, never the content read; secret paths are masked and command text is never written.

**11. Conformance suite, local and on demand.**
`scripts/harness-test.sh` runs a table of cases (tool, input, expected decision) against the core and against both adapters with synthetic JSON, checks the agent instruction files, validates the configuration JSONs with `jq empty`, and measures the average time per call. It is a specification-conformance check, not part of the regular test routine: it runs **on demand, locally**. There is no GitHub workflow for the harness (the user does not want harness gates on GitHub), and the existing `ci.yaml` and `build_and_release.yaml` do not call it; the suite checks that. It needs no network or secrets. `shellcheck` is an optional local lint for the scripts, not part of the suite. Each case carries the name of the specification scenario it proves (`# scenario: <requirement> / <scenario>`); the suite reads the scenario list from the specification (the change's delta spec, or `openspec/specs/agent-boundary/spec.md` after archive) and fails when a scenario has no case or a case names an unknown scenario, so the suite is the proof that every security and autonomy point works. The real behavior of the harnesses (the unconfirmed items) is verified by a manual checklist in a test repository, documented, because it needs `claude` and `agy` authenticated.

**12. Agent instructions: one common file, two thin harness files.**
`AGENTS.md` holds everything that applies to any agent: a short project overview (Go CLI with `cobra`; layers `cmd → internal/core → service → scraping/resource`), the verification routine (`go vet ./...`, `go test -race ./...`, `gofmt -l .`; the harness suite is **not** part of it and runs only when the user asks), test conventions (offline, using `internal/scraping/testdata/`; integration tests need the network and run only on demand; every new test must run in CI), commit and branch conventions (English conventional commits, branch + PR, push and merge belong to the user), a summary of the security boundaries (stay in the repository, never read secrets, no destructive commands, never work around a denied or asked command, ask the user instead) with a link to `docs/harness/local-guardrails.md`, the language convention, and the OpenSpec workflow (changes under `openspec/changes/`, `/opsx:*` commands, a task is checked only when its specified behavior is complete and verified). It stays short because it is always in context, and holds nothing the code already shows.
`CLAUDE.md` starts with `@AGENTS.md` (Claude Code does not read `AGENTS.md` itself) and adds only what is Claude-specific: what the messages starting with `guard:` mean (`ask` = user is consulted, `deny` = hard block) and where the log is; changes to `.claude/settings.json` and the hook take effect only after restarting the session (check with `/hooks` and `/permissions`); prefer `Read`/`Grep`/`Glob` over `cat`/`grep` in the shell; commit with several `-m` instead of `$(cat <<EOF)`, which asks; `.claude/settings.local.json` is personal and ignored; `.claude/skills` and `.claude/commands` are generated by OpenSpec and not edited by hand.
`GEMINI.md` imports `AGENTS.md` (`[VERIFY]` whether `agy` reads `AGENTS.md` natively and whether `@` imports work; if it reads it natively, the import line is dropped to avoid loading it twice) and adds only what is Antigravity-specific: the project hook in `.agents/hooks.json` answers `force_ask` or `deny`; permission settings are global and not versioned; never write to `~/.gemini`, since the template `.agents/agy-settings.example.json` is installed by the user; tool names (`run_command`, `view_file`, `write_to_file`…) and `/hooks`; `.agents/skills` and `.agents/workflows` are generated by OpenSpec and not edited by hand.
A rule lives in exactly one file: shared in `AGENTS.md`, otherwise in the harness file. None of the three contains secrets or machine-specific absolute paths. The suite checks that the three files exist, that `CLAUDE.md` imports `AGENTS.md`, and that none holds a home-directory path.

**13. Language convention.**
`README.md` is in Brazilian Portuguese; everything else is in English. The Portuguese already produced by this change (OpenSpec artifacts, `docs/harness/local-guardrails.md`, comments and messages in the scripts and hooks) is translated as part of this change. The guard's reason texts and the suite output are user-facing messages and follow the same rule. A language rule is also stated in `AGENTS.md`.

**14. File I/O only through the file tools.**
The agent reads with Read/Grep/Glob and writes with Edit/Write. The guard denies, with a message naming the right tool, the shell commands that duplicate them: writers (`sed` in any form, `tee`, `cp`, `mv`, `touch`, `install`, `ln`, `truncate`, `patch`, `rsync`, `git apply`, `curl -o`, `wget -O`), output redirects to a file (`>`, `>>`, `>|`, heredoc to a file) and readers (`cat`, `head`, `tail`, `less`, `more`, `bat`, `grep`/`egrep`/`fgrep`/`rg`, `find`, `ls`, `tree`, `awk`, `cut`, `jq`). Redirects to `/dev/null` and `2>&1` are allowed. Read-only `git`/`gh`, `go`, `gofmt` and a non-recursive `rm` are not affected. Mirrored as `deny` rules in `.claude/settings.json` and in the `agy` template; the hook is the primary enforcement. Edits and reads then always pass through the tools' own permission model (`acceptEdits`, path checks, protected files).

**15. One command per call.**
Chaining with `&&`, `||`, `;`, `&`, newline and `|` is `deny`, with a message asking for one command at a time. The splitter that already exists for nested destructive commands is reused, so a destructive command inside a chain is still reported as destructive. Literal `bash -c '...'` is still analyzed. `2>&1` and `/dev/null` redirects are not chaining.

**16. Scripting interpreters are denied.**
`python`, `python3`, `node`, `nodejs`, `perl`, `ruby`, `php`, `lua`, `deno` (inline, script file or heredoc) are `deny`: the project is Go, and an interpreter is the usual way around the file tools. This replaces the earlier `ask` for inline interpreters.

## Risks / Trade-offs

- [Without a sandbox, code the agent runs (`go test`, `go run`, scripts) reads whatever the user can] → Accepted and documented residual risk; keep the autonomous list short, ask on opaque commands, review diffs before running something on the host; a sandbox remains possible later.
- [A text guard can be bypassed: an assembled `$VAR`, `base64 | sh`, `python` reading a file] → Commands the parser cannot resolve become `ask`, not `allow`; deny `| sh`/`| bash`; document what is not covered.
- [Too many `ask`s on opaque commands reduce autonomy] → Measure in the suite which normal-cycle commands fall into `ask` and adjust the list.
- [Unconfirmed items in the Claude documentation: Edit/Write under the read key, Read outside the repository in `acceptEdits`, hook `ask` over `allow`, `cp` with an external target, `grep` naming a path under `deny`] → The manual verification task tests each one and the hook covers the gaps.
- [Antigravity hook fields (`toolCall.args`), tool names, the command path in `hooks.json` and the meaning of empty output may differ] → Adapter isolated from the core; verify on first run and adjust the adapter.
- [Antigravity configuration is global and outside the repository] → Versioned template and project hook; ask the user to install it and check with `/hooks`.
- [`agy` may not read `AGENTS.md` or may not accept imports in `GEMINI.md`] → Verify in the real CLI; fall back to the opposite arrangement (native read, or a self-contained `GEMINI.md`) and record it in the documentation.
- [The user has `blockReadsOutsideWorkingDirectories` on globally, so reads outside the repository are refused, not asked] → Document it; it is a safer choice than the request.
- [Denying `~/.claude` (per the spec) also blocks Claude Code's own memory and persisted tool results] → Documented; relaxing it is a spec, core and `settings.json` change.
- [Denying shell readers, chaining and pipes removes `go test | tail`, `git log | head` and `cd x && cmd`] → Accepted by the user; use `git log -n 20`, `go test -run`, and separate calls. The messages name the alternative.
- [The agent can still write through programs it is allowed to run: `go run`, `make`, an existing script] → Documented residual risk; those runs are not in the autonomous list and are asked in Claude Code.
- [Always-loaded instruction files drift from the code or grow] → Keep them short, one rule in one file, and re-check them when `dev-command-surface` lands.

## Migration Plan

1. Add configuration, guard, suite, documentation and agent instruction files, without changing the application code.
2. Run the local suite and the manual checklist for the unconfirmed items; adjust rules.
3. Open the PR from `chore/local-harness-guardrails`. Pushing and merging are the user's steps.
4. Rollback: revert the PR; the agent goes back to asking permission for everything and loses the guard and the instruction files. The `agy` template the user installed in `~/.gemini` must be removed by them.

## Open Questions

- Does the user want to turn off `blockReadsOutsideWorkingDirectories` in their global configuration to get asking instead of refusing? It is their adjustment and does not change this implementation.
- Turn on the native sandbox later as a second layer? A future decision, with no impact on this change's tasks.
- Add `Write all artifacts in English except README.md` to `openspec/config.yaml` so future OpenSpec artifacts follow the language convention automatically? Not part of this change; the user decides.

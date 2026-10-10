# Local guardrails for coding agents

Limits on what Claude Code and Antigravity (`agy`) may do in this repository. The agents run
directly on the host, with no container and no sandbox. The only boundary is the repository
directory: inside it the agent works without asking; anything that leaves it goes through the
user; secrets and destructive commands are denied.

Specification: `openspec/changes/local-harness-guardrails/` (archived under
`openspec/changes/archive/` once the change is closed).

## Threat model

The risk comes from the host and the user's credentials, not from the repository (which has no
database and no production variables). An agent has already used a wrong `.env` and run a
*delete* against a production database. Scenarios covered:

| Scenario | Example |
| --- | --- |
| Reading or copying secrets | `cat .env`, `~/.ssh/id_ed25519`, `~/.gemini/oauth_creds.json` |
| Leaving the repository | `/etc/hosts`, `~/x`, `../../other-project`, a symbolic link pointing outside |
| Destructive or out-of-scope command | `rm -rf`, `sudo`, `git push --force`, `psql`, `terraform apply`, `curl x \| sh` |
| Command that hides its target | `$(...)`, backtick, `bash -c`, `python3 -c`, `eval` |
| Weakening the guard itself | editing `.claude/`, `.agents/`, `scripts/guard*.sh`, `Makefile`, workflows, `AGENTS.md`, `CLAUDE.md`, `GEMINI.md` |
| Bypassing the file tools | `sed -i`, `cat > f`, `echo x >> f`, `tee`, `cp`, `python3 script.py` writing files |
| Packing many actions into one call | `a && b; c \| d` |

## Layers

1. **Permission rules** (`.claude/settings.json`): `acceptEdits` mode, a short list of autonomous
   commands (Go cycle and git/CI reads), `ask` and `deny`. `bypassPermissions` is disabled.
2. **Guard hook** (`PreToolUse`): inspects Bash and the file tools and decides
   `allow`/`ask`/`deny`. The logic lives in `scripts/guard-core.sh`, shared by the adapters
   `.claude/hooks/guard.sh` (Claude Code) and `.agents/hooks/guard-agy.sh` (Antigravity).
3. **Conformance suite** (`scripts/harness-test.sh`): feeds the core and both adapters with
   forbidden, ambiguous and allowed commands and fails if any decision differs. Every case names
   the specification scenario it proves, and the suite fails if a scenario has no case, so a passing
   run shows that every security and autonomy point works. It is **on demand**: it is not part of
   the regular test routine, and there is no GitHub workflow for it. Run it locally.

The hook **never returns `allow`**: when it has no objection it prints nothing and the decision is
left to the harness rules. Text it cannot parse becomes `ask`, never `allow`.

### What each layer covers and does not cover

| | Covers | Does not cover |
| --- | --- | --- |
| `settings.json` rules | Command prefixes, paths in the file tools, `cat`/`head`/`tail`/`sed` and redirects that name the file | `bash -c`, an absolute binary path, variables, scripts that open the file without naming it |
| Guard hook | The real command name (ignoring quotes, wrappers, binary path), chaining (`&&`, `\|\|`, `;`, `\|`, `&`, newlines), a literal `bash -c '...'`, paths through `realpath -m` (`..`, `~`, `$HOME`, symbolic links), secrets, writes to harness configuration | Code run by programs the agent launches, a `$VAR` assembled at run time, a two-step `base64 \| sh`, a glob that matches a symbolic link pointing outside |
| Suite | Regression of each rule, scenario coverage of the specification, parity between the adapters, timing, valid JSON, instruction files, language convention, on-demand wiring | Real harness behavior (see the checklist) |

### Residual risk (no sandbox)

No layer is a kernel boundary. A program the agent runs (`go test`, `go run`, `make`, a script)
reads whatever the user can read. The layers reduce the chance of the agent getting there unseen;
they do not remove it. In particular, the agent can still write files through programs it is
allowed to run (`go run`, `make`, an existing script). Mitigations: a short autonomous list, `ask`
on opaque commands, review the diff before running something new on the host. The native sandbox (`sandbox.enabled`) can be
turned on later as a second layer without redoing this.

## Decision table

| Input | Decision |
| --- | --- |
| Edit, read or search inside the repository with the file tools; a single `go build/test/vet/fmt/list`, `go mod tidy/download`, `gofmt`, `git status/diff/log/show/branch/switch/add/commit/stash`, `gh run list/view/watch`, `gh pr checks`; `2>&1` and `/dev/null` redirects | `allow` (no objection from the guard) |
| Shell writers: `sed` (any form), `tee`, `cp`, `mv`, `touch`, `install`, `ln`, `truncate`, `patch`, `rsync`, `tar`, `zip`, `unzip`, `git apply`, `curl -o`, `wget -O`; redirects to a file (`>`, `>>`, `>\|`, heredoc) | `deny` (use Edit or Write) |
| Shell readers: `cat`, `tac`, `head`, `tail`, `less`, `more`, `bat`, `grep`/`egrep`/`fgrep`/`rg`/`ag`, `find`/`fd`, `ls`, `tree`, `awk`, `cut`, `jq`, `sort`, `uniq`, `wc`, `diff`, `stat`, `file`, `git grep`; `< file` | `deny` (use Read, Grep or Glob) |
| Interpreters: `python*`, `node`, `perl`, `ruby`, `php`, `lua`, `deno`, `irb`, `pwsh`... (inline, script or heredoc) | `deny` |
| Chaining: `&&`, `\|\|`, `;`, `&`, newline, `\|` | `deny` (one command per call) |
| Absolute path, `~`, `$HOME`, a `..` that escapes, or a symbolic link pointing outside | `ask` |
| `$(...)`, backtick, `<(...)`, heredoc, `eval`, `bash -c`, `sh -c`, unresolved variable | `ask` |
| `git push` (not forced), `go get`, `go install`, `npm install`, any `gh` that is not a CI query | `ask` |
| Editing `.claude/`, `.agents/`, `.github/workflows/`, `Makefile`, `scripts/guard*.sh`, `scripts/harness-test.sh`, `AGENTS.md`, `CLAUDE.md`, `GEMINI.md` (file tool or shell) | `ask` |
| `.env*`, `*.pem`, `*.key`, `secrets/`, `~/.ssh`, `~/.aws`, `~/.config/gcloud`, `~/.kube`, `~/.gemini`, `~/.claude` | `deny` |
| Recursive/forced `rm`, `sudo`, `su`, `doas`, `chown`, `dd`, `mkfs*`, `shred` | `deny` |
| `git push --force/-f/--delete/+ref`, `git reset --hard`, `git clean -f`, `git filter-repo/branch` | `deny` |
| `psql`, `mysql`, `mongosh`, `redis-cli` and other database clients | `deny` |
| `terraform/tofu apply/destroy`, `kubectl delete/drain`, `gh repo delete`, `gh secret`, `gh auth token`, `npm/cargo publish` | `deny` |
| Anything `\| sh`/`\| bash`, `curl`/`wget` with the output piped to a shell or interpreter | `deny` |

`deny` holds even when nested (`cd /tmp && git push --force`, `bash -c "rm -rf x"`,
`find -exec rm -rf`). If a command matches both `deny` and `ask`, `deny` wins. Quoted text is
data, not a command: `git commit -m "mention rm -rf"` passes.

Effects worth knowing:

- Pipes and chains are gone: use `git log -n 20` instead of `git log | head`, `go test -run X`
  instead of `go test | grep X`, and separate calls instead of `cd x && cmd`. The denial message
  names the alternative. Quoted separators are data: `git commit -m "a && b"` passes.
- Chaining and the other shell rules are enforced by the hook. The `deny` entries in
  `settings.json` and in the `agy` template only mirror the command names; neither harness can
  express "no chaining" reliably in its rule syntax.
- `.env.example` is also denied (`.env*`). To allow it, adjust `_g_secret_path` in
  `scripts/guard-core.sh` and the `Read(**/.env*)` rules in `settings.json`.
- `git commit -m "$(cat <<'EOF' ... EOF)"` falls into `ask` (command substitution). Use several `-m`.
- Denying `~/.claude` includes Claude Code's own automatic memory and persisted tool results
  (`~/.claude/projects/...`), following the specification. If that gets in the way, relax the
  `~/.claude` rule in the specification, the core and `settings.json`.
- `blockReadsOutsideWorkingDirectories` in the user's own Claude Code configuration makes the file
  tools refuse reads outside the repository instead of asking. The project does not turn it on;
  turn it off in your own configuration if you prefer being asked.

## Running and extending

Prerequisites: `bash` and `jq`.

```sh
scripts/harness-test.sh                     # on demand: exits non-zero if any decision or scenario check fails
HARNESS_MAX_MS=50 scripts/harness-test.sh   # limit for the average time per call (default 100 ms)
bash scripts/guard-core.sh Bash "$PWD" 'git push --force'   # direct query: "deny<TAB>reason"
```

The suite is not part of the regular test routine (`go vet`, `go test -race`, CI on pull
requests) and no GitHub workflow runs it; run it when you change the guard, the rules or the specification, or when you want
proof that every point works.

To extend a list: edit `_g_rules` (commands), `_g_secret_path` (secrets) or `_g_protected`
(harness configuration and instruction files) in `scripts/guard-core.sh`; add the matching case to
`scripts/harness-test.sh` under a `scenario "<Requirement> / <Scenario>"` line (use
`case_ TOOL INPUT allow|ask|deny`) and run the suite. Keep the equivalent rules in
`.claude/settings.json` and `.agents/agy-settings.example.json`.

Scenario coverage: the suite reads the scenario names from the specification (the change's delta
spec, or `openspec/specs/agent-boundary/spec.md` after archive; override with `HARNESS_SPEC_FILE`).
It fails with `FAIL [coverage] missing scenario: ...` when a scenario has no case and with
`unknown scenario: ...` when a case names a scenario that does not exist, so a new scenario
cannot be added without a test that proves it.

## Log

Each `ask` or `deny` decision becomes a line in `.claude/logs/guard.log` (git-ignored):

```
2026-10-09T21:16:13-0300 Read deny secret access (path masked)
```

Format: `date tool decision reason`. The command text and the path of secrets are never written.
Only the hooks write the log (a direct query to the core and the suite do not write to the
repository); the file appears on the first `ask` or `deny` of a session.

```sh
tail -n 20 .claude/logs/guard.log
grep ' deny ' .claude/logs/guard.log
```

## Agent instruction files

| File | Content |
| --- | --- |
| `AGENTS.md` | Everything that applies to any agent: project overview, verification routine, test and commit conventions, security boundaries, language convention, OpenSpec workflow |
| `CLAUDE.md` | Starts with `@AGENTS.md` (Claude Code does not read `AGENTS.md` by itself) and adds only what is Claude-specific |
| `GEMINI.md` | Only what is Antigravity-specific. It does not import `AGENTS.md`: `agy` reads `AGENTS.md` natively and loads it together with `GEMINI.md`, and a `@` import in `GEMINI.md` did not resolve when tested |

A rule lives in exactly one file: shared in `AGENTS.md`, otherwise in the harness file. None of them
contains secrets or machine-specific paths. Editing any of them asks the user, like the rest of the
harness configuration, and the suite checks that they exist, that `CLAUDE.md` imports `AGENTS.md`,
that `GEMINI.md` does not, and that none holds a home-directory path (and that nothing but
`README.md` contains Portuguese).

## Language convention

`README.md` is written in Brazilian Portuguese. Everything else (documentation, OpenSpec artifacts,
instruction files, comments and messages of scripts and hooks, including the guard's reason texts
and the suite output) is written in English.

## Antigravity (`agy`)

The CLI only has **global** permission settings (`~/.gemini/antigravity-cli/settings.json`), so they
cannot be versioned. The repository provides:

- `.agents/hooks.json` + `.agents/hooks/guard-agy.sh`: a project hook with the same core logic. It
  asks with `force_ask` (ignores permissions already granted) and denies with `deny`.
- `.agents/agy-settings.example.json`: a template with `allowNonWorkspaceAccess: "off"`,
  `toolPermission: "request-review"` and `command(...)`, `read_file(...)`, `write_file(...)` rules.

The agent does not write to `~/.gemini`. To install the template, run this yourself from the
repository root (it merges without deleting your rules and keeps your other keys; the template's
`toolPermission` and `allowNonWorkspaceAccess` keys replace yours):

```sh
S=~/.gemini/antigravity-cli/settings.json
cp "$S" "$S.bak"
jq -s '
  .[0] as $cur | .[1] as $new
  | ($cur + ($new | del(.permissions)))
  | .permissions = (
      reduce ("allow","ask","deny") as $k ({};
        .[$k] = ((($cur.permissions[$k] // []) + ($new.permissions[$k] // [])) | unique))
    )' "$S.bak" .agents/agy-settings.example.json > "$S.new" && mv "$S.new" "$S"
```

To undo: `mv "$S.bak" "$S"`. Then check with `/hooks` and `/permissions` in `agy`.

Antigravity gaps and compensations:

| Gap | Compensation |
| --- | --- |
| No path glob in the rules | The hook resolves the path and asks/denies |
| The user's global rules add to the project's | Deny beats ask and allow; check `/permissions` |
| A hook `allow` may be ignored in `request-review` | The hook does not emit `allow`; only `force_ask` and `deny` |
| Tool names and `args` fields come from the documentation | See the checklist |

## Manual checklist (needs `claude` and `agy` authenticated)

Run it in a test repository, with `/permissions` and the debug log (`claude --debug`), after any
change to the hook, the rules or the harness versions. Record outcomes in the change's tasks or in
the pull request, not in this document.

| Item | How to verify |
| --- | --- |
| The hook runs and denies | Ask the agent for `rm -rf x`; expect a block with `guard: action denied` |
| Shell reads, writes and chains are denied | Ask the agent to `cat README.md`, to `sed -i` a file and to run `go vet ./... && go test ./...`; expect denials that name the right tool or ask for one command per call |
| Interpreters are denied | Ask the agent to write a file with a `python3` script; expect a denial |
| Edit/Write under `blockReadsOutsideWorkingDirectories` | Edit a file outside the repository with the key on |
| Read outside the repository in `acceptEdits` | Ask for `Read /etc/hosts` |
| A hook `ask` over an `allow` rule | Run `git push` with `Bash(git push *)` outside the `ask` list |
| `cp` with an external target | `cp README.md /tmp/x` |
| `grep` with a quoted path under `deny` | `grep x .env` |
| `ask` of `Edit(/.claude/**)` in `acceptEdits` | Edit `.claude/hooks/guard.sh` |
| Instruction files are loaded | `/memory` in Claude Code lists `CLAUDE.md` and the imported `AGENTS.md`; in `agy`, ask for a rule that only `AGENTS.md` states |
| Antigravity: tool names and `args` fields | `/hooks` in `agy`; call `run_command` and `view_file` |
| Antigravity: the `command` path in `hooks.json` and empty output meaning "no objection" | Run `agy` in the repository |

If a rule alone is not enough, the hook already covers the case with its own `ask`; adjust it and
record the result.

## Reverting

Revert the PR: the agent goes back to asking permission for everything and loses the guard and the
instruction files. Remove yourself what you installed in `~/.gemini` (`mv "$S.bak" "$S"`).

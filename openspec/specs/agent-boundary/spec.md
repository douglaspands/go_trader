# agent-boundary Specification

## Purpose
Define the limits on what coding agents may do on the host: autonomy inside the repository, user confirmation for any access outside it, secret protection, blocking of destructive commands, file reads and writes only through the file tools, one command per shell call, no scripting interpreters, shared and harness-specific agent instructions, a repository language convention, and on-demand verification that every one of these rules works in both harnesses.

## Requirements

### Requirement: Autonomy inside the repository
The agent SHALL run without a permission prompt: reading, searching and editing files inside the repository through the agent's file tools, except secrets, and the Go development-cycle commands: build, test, vet, formatting, module listing and read-only git commands, each as a single command.

#### Scenario: Editing a project file
- **WHEN** the agent edits a source file inside the repository with the file tools
- **THEN** the edit is applied without a permission prompt

#### Scenario: Reading a project file
- **WHEN** the agent reads or searches a file inside the repository with the file tools
- **THEN** the read happens without a permission prompt

#### Scenario: Running the project tests
- **WHEN** the agent runs `go test -race ./...`
- **THEN** the command runs without a permission prompt

### Requirement: Access outside the repository requires confirmation
Any access to a path outside the repository, by a file tool or by a shell command, SHALL require explicit user confirmation and SHALL NOT be approved automatically. This includes absolute paths, `~`, `..` that escapes the repository, and symbolic links that point outside it.

#### Scenario: Reading outside the repository
- **WHEN** the agent tries to read `/etc/hosts` or a file in the user's home directory
- **THEN** the user is consulted and the read happens only if they approve

#### Scenario: Escape through a relative path
- **WHEN** the agent runs a command with `../../other-project/file`, which resolves outside the repository
- **THEN** the user is consulted before the command runs

#### Scenario: Symbolic link
- **WHEN** the agent accesses a path inside the repository that is a symbolic link to somewhere outside it
- **THEN** the user is consulted before the access

### Requirement: Opaque commands require confirmation
A shell command whose target cannot be determined by reading its text, such as command substitution, `eval`, `bash -c` or `sh -c`, SHALL require user confirmation.

#### Scenario: Command substitution
- **WHEN** the agent runs a command that contains `$(...)` or a backtick
- **THEN** the user is consulted before the command runs

#### Scenario: Shell with inline code
- **WHEN** the agent runs `bash -c "..."` or `eval ...`
- **THEN** the user is consulted before the command runs

### Requirement: Files are read and written only through the file tools
The agent SHALL read files only with its search and read tools (Read, Grep, Glob) and write files only with its edit tools (Edit, Write). A shell command that reads or writes a file SHALL be denied, and the denial message SHALL name the tool to use instead. This covers file-writing commands (`sed`, `tee`, `cp`, `mv`, `touch`, `install`, `ln`, `truncate`, `patch`, `rsync`, `git apply`, `curl -o`, `wget -O`), output redirection to a file (`>`, `>>`, `>|`) and heredocs written to a file, and file-reading commands (`cat`, `head`, `tail`, `less`, `more`, `bat`, `grep`, `egrep`, `fgrep`, `rg`, `find`, `ls`, `tree`, `awk`, `cut`, `jq`). Redirects to `/dev/null` and file-descriptor redirects such as `2>&1` are not file access. Read-only `git` and `gh` CI queries, `go` and `gofmt` commands, and a non-recursive `rm` are not affected. The attempt SHALL be logged.

#### Scenario: Writing with sed
- **WHEN** the agent runs `sed -i s/a/b/ internal/core/app.go`
- **THEN** the command is denied and the message points to the Edit tool

#### Scenario: Writing with cat or echo
- **WHEN** the agent runs `cat > notes.md` with a heredoc, or `echo x >> notes.md`
- **THEN** the command is denied and the message points to the Write tool

#### Scenario: Similar writing commands
- **WHEN** the agent runs `tee out.txt`, `cp a.go b.go`, `touch x` or `patch -p1`
- **THEN** the command is denied

#### Scenario: Reading with a shell command
- **WHEN** the agent runs `cat README.md`, `grep -r TODO .`, `ls -la` or `find . -name x`
- **THEN** the command is denied and the message points to the Read, Grep or Glob tool

#### Scenario: Harmless redirects
- **WHEN** the agent runs `go test ./... 2>&1` or a command that redirects to `/dev/null`
- **THEN** the redirect does not make the command denied

### Requirement: Interpreters are denied
The agent SHALL NOT run a scripting-language interpreter (`python`, `python3`, `node`, `perl`, `ruby`, `php`, `lua`, `deno`), with inline code or with a script file, because it can write or read files without going through the file tools. The attempt SHALL be denied and logged.

#### Scenario: Inline interpreter
- **WHEN** the agent runs `python3 -c "..."` or `node -e "..."`
- **THEN** the command is denied

#### Scenario: Script file or heredoc
- **WHEN** the agent runs `python3 script.py` or `python3 <<EOF ... EOF`
- **THEN** the command is denied

### Requirement: One command per call
The agent SHALL run one command per shell call. Chaining with `&&`, `||`, `;`, `&`, a newline or a pipe `|` SHALL be denied, and the denial message SHALL tell the agent to run the commands one at a time. A literal `bash -c '...'` is still analyzed, and a destructive command nested in a chain SHALL still be reported as destructive.

#### Scenario: Chained commands
- **WHEN** the agent runs `go vet ./... && go test ./...`
- **THEN** the command is denied and the message asks for one command per call

#### Scenario: Pipe
- **WHEN** the agent runs `go test ./... | tail -20`
- **THEN** the command is denied

#### Scenario: Single command
- **WHEN** the agent runs `go vet ./...`
- **THEN** the command is not affected by this requirement

### Requirement: Secrets are never read
The agent SHALL NOT read, list the contents of, or copy `.env*` files, `*.pem`, `*.key`, the `secrets/` directory, nor `~/.ssh`, `~/.aws`, `~/.config/gcloud`, `~/.kube`, `~/.gemini` and `~/.claude`, by file tool or by shell. The attempt SHALL be denied without consulting the user and SHALL be logged.

#### Scenario: Reading .env
- **WHEN** the agent tries to read `.env` or `config/.env.production`, by file tool or by `cat`
- **THEN** the action is denied and the attempt is logged

#### Scenario: User credentials
- **WHEN** the agent tries to read `~/.ssh/id_ed25519` or `~/.gemini/oauth_creds.json`
- **THEN** the action is denied and the attempt is logged

### Requirement: Destructive commands are denied
The agent SHALL NOT run destructive or out-of-scope commands, such as forced recursive removal, privilege escalation, git history rewriting, database clients and package publishing. The attempt SHALL be denied and logged, including when nested in `&&`, `;`, `|` or `bash -c`.

#### Scenario: Recursive removal
- **WHEN** the agent runs `rm -rf build/`
- **THEN** the command is denied and logged

#### Scenario: Privilege escalation
- **WHEN** the agent runs `sudo`, `su`, `chown`, `dd`, `mkfs` or `shred`
- **THEN** the command is denied and logged

#### Scenario: Git history
- **WHEN** the agent runs `git push --force`, `git reset --hard`, `git clean -fdx`, `git filter-repo` or `git filter-branch`
- **THEN** the command is denied and logged

#### Scenario: Infrastructure and publishing
- **WHEN** the agent runs `terraform apply`, `kubectl delete`, `gh repo delete`, `gh secret` or `npm publish`
- **THEN** the command is denied and logged

#### Scenario: Downloaded code executed
- **WHEN** the agent runs `curl` or `wget` with the output piped to a shell
- **THEN** the command is denied and logged

#### Scenario: Nested command
- **WHEN** the agent runs `cd /tmp && git push --force`
- **THEN** the command is denied even though it is chained

#### Scenario: Database client
- **WHEN** the agent runs `psql` or `mongosh`
- **THEN** the command is denied

### Requirement: External and irreversible actions require confirmation
The agent SHALL ask for user confirmation before a non-forced `git push`, `go get`, `go install`, `npm install`, any `gh` command that is not a CI query, and before editing `.claude/`, `.agents/`, `.github/workflows/`, `Makefile`, `AGENTS.md`, `CLAUDE.md` or `GEMINI.md`.

#### Scenario: Push
- **WHEN** the agent runs `git push origin chore/local-harness-guardrails`
- **THEN** the user is consulted before the command runs

#### Scenario: Harness configuration
- **WHEN** the agent tries to edit `.claude/hooks/guard.sh`
- **THEN** the user is consulted, even in accept-edits mode

#### Scenario: Agent instructions
- **WHEN** the agent tries to edit `AGENTS.md`, `CLAUDE.md` or `GEMINI.md` with a file tool
- **THEN** the user is consulted, even in accept-edits mode

#### Scenario: Agent instructions through the shell
- **WHEN** the agent tries to change `AGENTS.md`, `CLAUDE.md` or `GEMINI.md` with a shell command such as a redirect
- **THEN** the command is denied, because shell writes are denied

### Requirement: Parity between Claude Code and Antigravity
For the same command or path, both harnesses SHALL produce the same decision (allow, ask or deny), through shared decision logic. Where a harness cannot enforce a rule, the gap SHALL be documented together with its compensation.

#### Scenario: Same decision
- **WHEN** the same destructive command is submitted to the guard through the Claude Code and Antigravity adapters
- **THEN** both return the deny decision

### Requirement: Decision logging
Every ask or deny decision of the guard SHALL be written to a local project log, with date, tool, decision and reason, without recording secret contents. The log SHALL stay out of version control.

#### Scenario: Denied attempt
- **WHEN** the guard denies a command
- **THEN** a line with date, tool, decision and reason is appended to the log

### Requirement: Guard performance
The guard SHALL decide in under 100 ms per call for a typical command, so it does not slow down the agent cycle.

#### Scenario: Typical command
- **WHEN** the guard evaluates `go test ./...`
- **THEN** the decision is returned in under 100 ms

### Requirement: Versioned agent instructions
The repository SHALL contain `AGENTS.md` with the instructions common to every agent, `CLAUDE.md` with only what is specific to Claude Code, and `GEMINI.md` with only what is specific to Antigravity. `CLAUDE.md` and `GEMINI.md` SHALL import `AGENTS.md` instead of repeating it. None of the three SHALL contain secrets, credentials or absolute paths specific to a machine.

#### Scenario: Common instructions
- **WHEN** an agent of either harness starts a session in the repository
- **THEN** it has the shared instructions from `AGENTS.md`: project overview, verification commands, test, commit and branch conventions, a summary of the security boundaries with a link to the full document, the language convention and the OpenSpec workflow

#### Scenario: Harness-specific instructions
- **WHEN** a rule applies to only one harness
- **THEN** it appears only in that harness's file (`CLAUDE.md` or `GEMINI.md`), and a rule that applies to both appears only in `AGENTS.md`

#### Scenario: Missing import
- **WHEN** `CLAUDE.md` does not import `AGENTS.md`
- **THEN** the regression suite fails and identifies the file

### Requirement: Language convention
`README.md` SHALL be written in Brazilian Portuguese. Every other repository artifact, including OpenSpec artifacts, documentation, agent instructions and the comments and messages of scripts and hooks, SHALL be written in English.

#### Scenario: New document
- **WHEN** a new document, script or instruction file is added outside the README
- **THEN** it is written in English

#### Scenario: README
- **WHEN** the README is changed
- **THEN** it stays in Brazilian Portuguese

### Requirement: Automatic verification of the rules
The repository SHALL have an executable suite that proves every security and autonomy point of this specification works. It submits forbidden, ambiguous and allowed commands to the guard and fails if any decision differs from the expected one, and checks that the agent instruction files exist. Every scenario of this specification SHALL be covered by at least one suite case that names it, and the suite SHALL fail when a scenario has no case or a case names a scenario that does not exist. The suite runs on demand, locally; it SHALL NOT run as part of the regular test executions and SHALL NOT be wired into any GitHub workflow. It SHALL need no secrets, network access or dependencies outside the repository.

#### Scenario: Regression of a rule
- **WHEN** a change to the guard stops denying `git push --force`
- **THEN** the suite fails and identifies the case

#### Scenario: Scenario without a case
- **WHEN** a scenario is added to this specification and no suite case names it
- **THEN** the suite fails and lists the uncovered scenario

#### Scenario: On demand only
- **WHEN** the regular test routine runs (`go vet ./...`, `go test -race ./...`, CI on a pull request)
- **THEN** the harness suite does not run, no GitHub workflow runs it, and it runs only when requested

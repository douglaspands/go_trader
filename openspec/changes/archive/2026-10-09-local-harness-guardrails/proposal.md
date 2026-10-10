# Proposal

## Why

Coding agents (Claude Code and Antigravity `agy`) will build the next features, including the React frontend, running directly on the host with no container or sandbox. The repository has no versioned agent configuration today, so every command raises a permission prompt and nothing stops an agent from reading or changing files outside the project. An agent has already used a wrong `.env` and run a *delete* against a production database. This change gives agents full autonomy inside the repository and requires explicit permission for anything that leaves it, with hard blocks on secrets and destructive commands. It also gives them the instructions they are missing: a shared `AGENTS.md`, plus harness-specific `CLAUDE.md` and `GEMINI.md`.

## What Changes

- Version `.claude/settings.json` with `acceptEdits` mode, a curated list of autonomous Go-cycle commands, and ask/deny rules.
- Create the `PreToolUse` guard hook, in bash and `jq` (fast), that decides `allow`/`ask`/`deny` for Bash and for the file tools: a path outside the repository asks, an opaque command asks, secrets and destructive commands are denied. Every `ask` or `deny` decision is written to a local log.
- Deliver the Antigravity equivalent: a project hook in `.agents/hooks.json` reusing the same logic, and a global `settings.json` template (`allowNonWorkspaceAccess` off and `command(...)` rules) that the user installs themselves because the CLI only has global configuration.
- Make the agent read and write files only through its file tools: shell readers and writers (`cat`, `sed`, `tee`, redirects, `grep`, `ls`, and similar) are denied with a message naming the right tool, scripting interpreters (`python`, `node`, `perl`...) are denied, and chaining (`&&`, `;`, `|`...) is denied so the agent runs one command per call.
- Create `scripts/harness-test.sh`, a suite that feeds the guard forbidden, ambiguous and allowed commands, checks the expected decision, and proves that every security and autonomy scenario of the specification has a passing case. It runs on demand, locally; there is no GitHub workflow for it and it is not part of the regular test executions.
- Add versioned agent instructions: `AGENTS.md` with what is common to every agent, `CLAUDE.md` with only what is specific to Claude Code, and `GEMINI.md` with only what is specific to Antigravity. The last two import `AGENTS.md` instead of repeating it. Editing these three files asks for confirmation, like the rest of the harness configuration.
- Adopt a language convention: `README.md` is written in Brazilian Portuguese; every other artifact (OpenSpec artifacts, docs, agent instructions, script comments and messages) is written in English. The Portuguese text already produced by this change is translated.
- Adjust the repository `.gitignore` for `.claude/settings.local.json` and the guard log.
- Document the threat model, known limits and step-by-step setup in `docs/harness/local-guardrails.md` and in the README.

No change to the application's Go code. No **BREAKING** change for users of `trader`.

## Capabilities

### New Capabilities
- `agent-boundary`: limits on what coding agents may do on the host: autonomy inside the repository, confirmation to leave it, secret protection, blocking of destructive commands, shared and harness-specific agent instructions, the repository language convention, and automatic verification of these rules.

### Modified Capabilities

## Impact

- New files: `.claude/settings.json`, `.claude/hooks/guard.sh`, `.agents/hooks.json`, `.agents/agy-settings.example.json`, `scripts/harness-test.sh`, `docs/harness/local-guardrails.md`, `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`.
- Changed: `.gitignore` and `README.md`.
- Host dependencies: `bash` and `jq` (already installed).
- Antigravity depends on a manual user step in `~/.gemini/antigravity-cli/settings.json`, outside the repository.
- Does not change `go.mod`, the application code or the existing CI and release workflows.
- Replaces the earlier container idea: container isolation is out of scope by the user's decision.

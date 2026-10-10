# Proposal

## Why

The archived change `local-harness-guardrails` shipped the guard, the adapters and the
conformance suite, but left two verification tasks open because they need real, authenticated
clients: 3.3 (Claude Code) and 4.2 (Antigravity `agy`). The suite proves the decision logic with
synthetic JSON; it cannot prove that Claude Code and `agy` actually call the hooks, pass the fields
the adapters expect, and honor `ask` and `deny` the way the design assumes. Until this is observed,
the design's "unconfirmed" items remain risks, and `docs/harness/local-guardrails.md` has a manual
checklist nobody has completed.

## What Changes

- Run the manual checklist from `docs/harness/local-guardrails.md` in a trusted Claude Code test
  repository and record the observed result of each item, covering the two items still open from
  3.3: Edit/Write under the read key, and a hook `ask` over an `allow` rule.
- Run the Antigravity part of the checklist in the real `agy`: tool names and `args` field names
  (`/hooks`), the command path in `.agents/hooks.json`, and that empty output means "no objection".
- Adjust `.claude/hooks/guard.sh`, `.agents/hooks/guard-agy.sh`, `.agents/hooks.json`,
  `.claude/settings.json` or `.agents/agy-settings.example.json` only where an observation shows the
  current behavior differs from the specification, with a suite case written first.
- Record the observed results in this change's tasks (they stay in the archive, not in the permanent
  documentation), and add the requirement that harness behavior is verified in the real clients.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `agent-boundary`: adds a requirement that the guard's behavior in Claude Code and Antigravity is
  verified in the real clients, with the observed results recorded. No existing requirement changes.

## Impact

- Files read or run: `docs/harness/local-guardrails.md` (checklist), `.claude/`, `.agents/`,
  `scripts/harness-test.sh`.
- Files that may change, only if an observation requires it: the two adapters, `.agents/hooks.json`,
  `.claude/settings.json`, `.agents/agy-settings.example.json`, `scripts/guard-core.sh`,
  `scripts/harness-test.sh` and the checklist in `docs/harness/local-guardrails.md`.
- No application Go code. No new dependencies.
- Needs the user: a trusted test repository for `claude`, and an authenticated `agy`. Editing the
  harness configuration asks the user as usual.

# Design

## Context

See `proposal.md` - Why. The archived change `local-harness-guardrails` left tasks 3.3 and 4.2
open. Already observed in real Claude Code sessions (recorded in the archived task 3.3): the hook
runs and blocks `rm -rf`; a Read outside the repository is logged as `ask`; `cp` to an external
target asks; `grep x .env` is denied; writes under `.agents/` and `.claude/` ask. Not yet observed:
Edit/Write under the read key, and a hook `ask` over an `allow` rule. Nothing has been observed in
the real `agy`.

Constraints:
- A headless `claude -p` run in an untrusted directory ignores the repository's `allow` rules, so
  it proves nothing. The Claude Code checks need a **trusted** directory, which the user accepts in
  an interactive session.
- `agy` needs the user's authentication, and its permission settings are global
  (`~/.gemini/antigravity-cli/settings.json`). The agent never writes to `~/.gemini` or `~/.claude`.
- The agent cannot start `claude` or `agy` against itself from inside a guarded session in a way
  that proves anything; the checks are run by the user in their own terminal, and the results are
  reported back.

## Goals / Non-Goals

**Goals:**
- A recorded result for every Claude Code and `agy` item of the checklist.
- The adapters and configuration match what the real clients send and accept.

**Non-Goals:**
- New guard rules or new specified behavior beyond the verification requirement.
- Turning on the native sandbox, or changing the user's global configuration.
- Automating the checks in CI (the user does not want harness gates on GitHub).

## Decisions

**1. The user runs the clients; the agent prepares and records.**
The agent writes the exact steps and expected outcomes for each item (from the checklist), the user
runs them in a trusted test copy and in `agy`, and reports what happened. The agent records each
result in `tasks.md` and fixes only what the observation requires. Alternative: drive `claude`
headless from the agent; rejected because untrusted headless runs ignore `allow` rules, and the
guard rightly stops a session from launching a nested one with the repository's own permissions.

**2. A test copy of the repository, not the working repository.**
The checks that need destructive-looking or outside-repository actions (Edit under the read key,
a hook `ask` over an `allow` rule) run in a disposable copy outside this repository, created by the
user, so a wrong result cannot touch real work. The copy is trusted once in an interactive session.

**3. Fix with the suite first.**
If an observation differs from the specification, a case is added to `scripts/harness-test.sh`
under the matching scenario and seen failing before the guard, adapter or configuration changes.
Where the hook is not enough, the rule is added to `settings.json` or the `agy` template; where a
rule alone is not enough, the hook covers it.

**4. Results stay in the tasks.**
Per the checklist's own rule, observed results are recorded in this change's tasks (kept in the
archive) and not in `docs/harness/local-guardrails.md`, which keeps only the procedure. The
document changes only if the procedure itself was wrong or incomplete.

## Risks / Trade-offs

- [The real client behaves differently from the documentation in a way the hook cannot compensate] →
  Record it, document the gap and its compensation in `local-guardrails.md` as a limitation, and
  surface it to the user instead of narrowing the specified behavior silently.
- [`agy` does not honor an empty output as "no objection", or uses other tool or `args` names] →
  Adjust the adapter and the `matcher` in `.agents/hooks.json`; the core is untouched.
- [The checks depend on the user's time and tools, so the change can stay open] → Each item is
  independent and recorded as soon as observed; the change is archived only when all are recorded or
  the user explicitly defers one.
- [Editing harness configuration asks the user each time] → Expected; it is the guard working.

# Spec Delta

## ADDED Requirements

### Requirement: Verification in the real clients
The guard's behavior SHALL be verified in the real Claude Code and Antigravity clients, not only with synthetic input, and the observed result of each checklist item SHALL be recorded in the tasks of the change that ran it. A difference between the observed behavior and this specification SHALL be fixed in the guard, the adapters or the harness configuration, with a suite case that fails before the fix.

#### Scenario: Claude Code checklist
- **WHEN** the manual checklist is run in a trusted test repository with Claude Code
- **THEN** each Claude Code item has a recorded result, including Edit/Write under the read key and a hook `ask` over an `allow` rule

#### Scenario: Antigravity checklist
- **WHEN** the manual checklist is run in the real `agy`
- **THEN** the tool names, the `args` field names, the command path in `hooks.json` and the meaning of empty hook output each have a recorded result

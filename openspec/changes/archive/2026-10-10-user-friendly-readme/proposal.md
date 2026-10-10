# Proposal

## Why

The current `README.md` is technical and developer-focused from the opening lines, jumping immediately into terminal commands, shell exit codes, and low-level flags. For retail investors and non-developers who want to calculate monthly investment allocations across B3 stocks and REITs (FIIs), this presents a barrier to entry and fails to communicate the motivation, benefits, safety, and target audience of the tool.

## What Changes

- Rewrite `README.md` with a clear, user-centric structure in Brazilian Portuguese following README best practices.
- Add an engaging project overview highlighting motivation, value proposition, and target audience (investors doing monthly B3 allocation, privacy-minded users, and automation enthusiasts).
- Add a beginner-friendly quickstart guide explaining how to download pre-built binaries for Windows, macOS, and Linux and run them step-by-step.
- Reorganize practical CLI examples around real-world investment scenarios (single quote, multi-asset quotes, stock balance, REIT balance, combined security balance, and CSV export).
- Add an FAQ addressing common non-technical questions (data provider origin, execution safety without broker credentials, internet connectivity, error handling).
- Retain all technical developer instructions (architecture, compilation, unit/integration testing, exit codes, CSV formula sanitization, and AI agent guardrails) in a dedicated "Para Desenvolvedores" section at the bottom.

## Capabilities

### New Capabilities

*(None. This change does not alter system behavior or CLI specifications; it updates user-facing documentation only.)*

### Modified Capabilities

*(None. All CLI flags, outputs, exit codes, and domain logic remain unchanged.)*

## Impact

- Documentation: `README.md` is completely restructured and rewritten.
- Code & APIs: None. No source code or tests are modified.
- Dependencies: None.

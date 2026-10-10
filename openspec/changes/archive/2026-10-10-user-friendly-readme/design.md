# Design

## Context

The repository provides a Go CLI (`trader`) that scrapes StatusInvest and calculates monthly portfolio purchase balances for B3 stocks and REITs (FIIs). Per `AGENTS.md`, `README.md` is maintained in Brazilian Portuguese, while internal documentation, specs, and agent instructions are in English.

The existing `README.md` lists CLI commands and terminal outputs alongside low-level developer details (exit codes, `stdout`/`stderr` channels, build commands, and AI agent harness instructions). To make the project accessible to retail investors and non-developers, the document needs a user-first visual hierarchy and storytelling approach while preserving all developer information.

## Goals / Non-Goals

**Goals:**
- Present a clear value proposition, motivation, and target audience upfront.
- Structure a beginner-friendly quickstart (downloading binaries, opening PowerShell/Terminal, running commands).
- Frame practical examples around realistic investor workflows (e.g., allocating a monthly budget of R$ 1.000 across multiple assets).
- Add an FAQ addressing safety, privacy, data sources, and error handling.
- Preserve all existing technical guidelines (build targets, unit/integration tests, exit codes, CSV formula sanitization, AI agent guardrails) in a dedicated developer section.

**Non-Goals:**
- Modifying CLI arguments, flags, output formats, or domain logic in Go code.
- Moving developer docs to separate external files (the README remains the single consolidated hub).
- Modifying other repository docs or agent instruction files (`AGENTS.md`, `GEMINI.md`, `CLAUDE.md`).

## Decisions

### 1. Two-Tier Document Structure
- **Decision**: Organize the README into an End-User Tier (Header -> Motivation -> Target Audience -> Features -> Quickstart -> Practical Examples -> FAQ) followed by a Developer Tier (`## 🛠️ Para Desenvolvedores`).
- **Rationale**: Non-developers get immediate clarity and actionable usage examples without being overwhelmed by build tools and exit codes. Developers still have a direct reference at the bottom.
- **Alternatives Considered**: Splitting into `README.md` and a separate `DEVELOPMENT.md`. Rejected to keep repository root self-contained and preserve existing workflow references.

### 2. Tone and Terminology
- **Decision**: Use clear, conversational Brazilian Portuguese for the user-facing sections, explaining concepts like tickers, monthly investment allocations (*aportes mensais*), and remainder amounts (*saldo restante*).
- **Rationale**: Bridges the gap between finance terms and CLI commands.

### 3. Reassurance & Security Messaging
- **Decision**: Add an FAQ explicitly clarifying that GoTrader does not require broker credentials, does not execute transactions on financial exchanges, and processes data entirely on the user's machine.
- **Rationale**: Financial tools often create apprehension; clarity builds trust.

## Risks / Trade-offs

- **[Risk]** Developers could overlook build or testing commands if placed lower in the document.  
  → **Mitigation**: Use a clear section title (`## 🛠️ Para Desenvolvedores`) and structured subheadings (`Compilação local`, `Testes e Qualidade`, `Desenvolvimento com Agentes de IA`).
- **[Risk]** Examples becoming inconsistent with actual CLI specifications.  
  → **Mitigation**: Verify all flags (`--amount`, `--stocks`, `--reits`, `--csv`, `--no-color`) and sample outputs against `openspec/specs/cli/spec.md` and `openspec/specs/purchase-balance/spec.md`.

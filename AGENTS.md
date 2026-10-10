# AGENTS.md

Instructions shared by every coding agent working in this repository. Harness-specific
notes live in `CLAUDE.md` (Claude Code) and `GEMINI.md` (Antigravity); nothing here is
repeated there.

## Project

`trader` is a Go CLI (`cobra`) that scrapes stock and REIT quotes from StatusInvest and
computes purchase balances. Layers: `main.go` starts `internal/core`, which wires the
`cmd` commands to `internal/service`, which uses `internal/scraping` (HTTP + HTML parsing)
and `internal/resource` (domain types). The Go version is declared in `go.mod`.
Behavior is specified under `openspec/specs/`.

## Verify before declaring work done

```sh
go vet ./...
go test -race ./...
gofmt -l .        # must print nothing
```

`make test/unit` also writes `coverage.out`, which is git-ignored.

The harness suite (`scripts/harness-test.sh`) is not part of this routine: run it on demand
only, when the user asks.

## Tests

- Unit tests are offline. Scraping tests read HTML from `internal/scraping/testdata/`.
- Integration tests (build tag `integration`, `make test/integration`) hit the real site;
  run them only on request.
- Every new test must run in CI (`.github/workflows/ci.yaml` runs `go test -race ./...`).

## Git

- Conventional commits in English (`feat:`, `fix:`, `chore:`, `docs:`).
- Work on a branch and open a PR against `main`. Pushing and merging are the user's steps.

## Security boundaries

- Stay inside the repository. Anything outside it needs the user's approval.
- Never read or write secrets (`.env*`, keys, `~/.ssh`, `~/.aws`, `~/.gemini`, `~/.claude`).
- No destructive or out-of-scope commands (`rm -rf`, forced git operations, database
  clients, publishing, infrastructure changes).
- Read files with Read, Grep and Glob and write them with Edit and Write. Shell readers and
  writers (`cat`, `grep`, `ls`, `find`, `sed`, `tee`, `cp`, redirects) and scripting
  interpreters (`python`, `node`, `perl`...) are denied.
- Run one command per shell call: no `&&`, `||`, `;`, `&` or `|`.
- Editing the harness configuration or these instruction files asks the user first.
- When a command is denied or asked, do not work around it; ask the user instead.

Details, limits and the decision table: `docs/harness/local-guardrails.md`.

## Language

`README.md` is in Brazilian Portuguese. Everything else (docs, OpenSpec artifacts,
instruction files, comments, script messages) is in English.

## OpenSpec workflow

Changes live under `openspec/changes/<name>/` and are driven with the `/opsx:*` commands
(`propose`, `apply`, `archive`...). Mark a task `- [x]` only when its specified behavior
is fully implemented and verified. Do not narrow or defer specified behavior silently;
surface it and ask.

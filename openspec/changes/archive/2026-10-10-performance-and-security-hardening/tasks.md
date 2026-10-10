# Tasks

## 1. Input and output hardening

- [x] 1.1 Add `validateAmount` in `cmd/common.go` (finite, > 0, <= 1e12) and call it first in the three `purchase-balance` `RunE`s; verify with `cmd` tests for `inf`, `NaN`, `0`, `-100` and `1e13` that assert stderr has an invalid-amount error, exit code 1, no service call recorded, and that the test finishes (no hang)
- [x] 1.2 Limit tickers to 12 characters in `scraping.tickerPattern`; verify with stock and REIT scraping tests that a 13-character ticker returns an invalid ticker error and the `httptest` server receives no request, and that a 12-character ticker is accepted
- [x] 1.3 Add `tools.CleanText` (removes runes where `unicode.IsControl`) and `tools.CSVSafe` (prefixes `'` to text starting with `=`, `+`, `-`, `@`); verify with table-driven tests in `internal/tools` covering ESC sequences, BEL, C1 controls, each formula prefix, and text that must stay unchanged
- [x] 1.4 Build provider-text cells in `cmd/{stock,reit,security}.go` through a `cell` helper that applies `CleanText`, plus `CSVSafe` when `--csv` is set; verify with `cmd` tests that a name containing `\x1b[2J` is printed without ESC and that `--csv` turns `=HYPERLINK("http://x")` into `'=HYPERLINK("http://x")`, while the plain table leaves it unprefixed

## 2. Shared fetcher

- [x] 2.1 Add `GetMaxConcurrentRequests()` (4) to `config.Config` and the test fakes that implement it; verify with a `config` test and `go vet ./...`
- [x] 2.2 Create `scraping.Fetcher` (shared `http.Client` with timeout from config, semaphore sized by config, test option for the transport) replacing `getHtml`; drop the manual `Accept-Encoding` and `Connection` headers and the gzip branch; verify the existing gzip and invalid-gzip scraping tests still pass against the new fetcher
- [x] 2.3 Cap the decompressed body at 10 MiB in the fetcher; verify with tests that a body of exactly 10 MiB is accepted and one byte more returns an error and no security, for stocks and REITs
- [x] 2.4 Add the `CheckRedirect` policy (https only, same host as the first request, at most 10 hops); verify with `httptest.NewTLSServer` tests that a same-host https redirect is followed, and that redirects to another host or to `http` end with an error
- [x] 2.5 Inject one fetcher into `NewStockScraping` and `NewReitScraping` from `core.setup`; verify `internal/core` tests pass and the test helpers build scrapers through the new constructors
- [x] 2.6 Remove the stock description extraction (design D10); verify the stock scraping tests still pass and `Security.Description` is empty
- [x] 2.7 Complete the 12 truncated entries in `scraping.userAgents` (found during manual testing: StatusInvest answers 403 to every User-Agent whose last version number ends in `.`, so about a third of the requests failed at random, which `main` also hid by dropping the tickers); verify with `TestUserAgentsAreComplete` and by running `security purchase-balance` with 8 stocks and 5 REITs against the real site, all 13 quoted with 4 requests in flight

## 3. Concurrent lists and visible failures

- [x] 3.1 Add `resource.TickerFailure` and `resource.QuoteFailuresError` (its message lists every ticker and reason); verify with a `resource` test of the error message
- [x] 3.2 Change scraping `List{Stocks,Reits}ByTickers` to fetch concurrently and return `([]*Security, []*TickerFailure)` with successes in input order; verify with tests that (a) a server answering in reverse order still yields input order, (b) a failing ticker is reported with its reason and the others are returned, and (c) an atomic in-flight counter never exceeds 4 for 10 tickers
- [x] 3.3 Propagate the new signatures through `service.{Stock,Reit}Service` and update the service and `cmd` fakes; verify `go build ./...` and the service list tests
- [x] 3.4 Make `PurchaseBalancesBySecurities` fetch stocks and REITs concurrently and return `(*PurchaseBalance, error)`, returning a `*QuoteFailuresError` without computing a balance when any ticker fails; verify with service tests for one failure, failures on both sides, and the success path keeping stocks before REITs, plus a test that 3 stocks + 3 REITs never exceed 4 requests in flight
- [x] 3.5 In `cmd`, `list` prints one `Warning: <ticker>: <reason>` line per failure to stderr with exit code 0, and the three `purchase-balance` commands return the quote-failure error (stdout empty, exit 1); verify with `cmd` tests for each scenario in `specs/cli/spec.md`, including stderr being empty when every ticker succeeds
- [x] 3.6 Update the error-behavior section of `README.md` (Portuguese) for invalid amounts, list warnings and `purchase-balance` aborting on quote failures; verify the documented commands and messages match the `cmd` test expectations

## 4. CI and release workflows

- [x] 4.1 In `build_and_release.yaml`, pass `TAG_NAME` only through `env:` and reference it as `"$TAG_NAME"` in every `run:`; verify that no `${{` expression remains inside any `run:` block
- [x] 4.2 Set top-level `permissions: contents: read` and grant `contents: write` only to the jobs that upload; verify by reading the workflow that the `test` job has read-only permissions
- [x] 4.3 Add `govulncheck` to the release `test` job and switch `ci.yaml` to `go-version-file: go.mod`; verify both workflows declare the same Go source and run `govulncheck`
- [x] 4.4 Add a `checksums` job (needs all build jobs) that downloads the release assets, writes `SHA256SUMS` and uploads it; verify the job's `needs`, permissions and commands by reading the workflow
- [x] 4.5 Pin `actions/checkout` and `actions/setup-go` to full commit SHAs (tag in a comment) and `govulncheck` to an exact version, resolving them with the user's approval since it needs network access; verify that no `@v<major>` or `@latest` reference remains in `.github/workflows/`

## 5. Integration check

- [x] 5.1 Run `go vet ./...`, `go test -race ./...` and `gofmt -l .`; verify all pass and `gofmt` prints nothing
- [x] 5.2 Run `openspec validate performance-and-security-hardening --strict`; verify it reports no errors

## Workflow follow-up

- On request, run `make test/integration` against the real site to confirm the concurrent fetch and redirect policy work with StatusInvest.
- Open the PR from `feat/performance-and-security-hardening` against `main`; pushing and merging are the user's steps.
- Archive the change after merge with `/opsx:archive`.

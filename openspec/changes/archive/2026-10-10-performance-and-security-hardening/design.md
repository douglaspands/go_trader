# Design

## Context

See proposal.md for the motivation. The current flow, for every command:

```
cmd --> service --> scraping.getHtml(url)  (new http.Client per call)
                     |
List*ByTickers: for each ticker, in sequence; errors dropped
PurchaseBalancesBySecurities: all stocks, then all REITs
```

- `scraping.getHtml` sets `Accept-Encoding: gzip, deflate` by hand and decodes only gzip. Setting the header by hand also turns off the transparent decompression of Go's `http.Transport`.
- `stockScraping` and `reitScraping` are built separately in `internal/core/app.go`, each one from `config`.
- `service.List*ByTickers` and `PurchaseBalancesBySecurities` return only successes. The `cmd` tests fake these interfaces (`cmd/mocks_test.go`).
- `--amount` is a plain `Float64` flag, so `pflag` accepts `inf`, `NaN` and negative values.
- No new dependency is allowed unless it is needed: the module has four direct dependencies, and `golang.org/x/sync` is not one of them.

## Goals / Non-Goals

**Goals:**
- One HTTP client and one concurrency limit shared by every request of a command run.
- Failures carried from the scraping layer up to `cmd` without losing which ticker failed or why.
- Each hardening rule lives in one place: the fetcher for network rules, `cmd` for input and output rules.

**Non-Goals:**
- `stock get` / `reit get` keep their current message (`ticker "<x>" not found!`). Showing the failure reason there is a separate change.
- No retry, backoff or caching of quotes.
- No configurable concurrency through a flag or environment variable. The limit is a fixed value in `config`.
- No translation of the existing Portuguese spec text (see Risks).

## Decisions

### D1. A shared `fetcher` in `scraping`, built once in `core`

```
core.setup()
  cfg     := config.NewConfig()
  fetcher := scraping.NewFetcher(cfg)          // http.Client + semaphore(4)
  stock   := scraping.NewStockScraping(cfg, fetcher)
  reit    := scraping.NewReitScraping(cfg, fetcher)
```

The fetcher owns the `*http.Client` (timeout from `config`, `CheckRedirect`, default transport) and a buffered channel of capacity `cfg.GetMaxConcurrentRequests()` (4) used as a semaphore. `Fetch(url) ([]byte, error)` takes a slot only around the HTTP exchange and the body read, so HTML parsing does not hold a slot.

Because both scrapers share one fetcher, the limit of 4 holds even when `security purchase-balance` fetches stocks and REITs at the same time.

- *Alternative:* a package-level client and semaphore. Rejected because tests could not get an isolated limit or swap the transport.
- *Alternative:* the limit inside each `List*` function. Rejected because two lists running together would allow 8 requests in flight.

For tests, an option lets the caller pass an `*http.Client` transport, so a test can use `httptest.NewTLSServer` and its trusted transport to exercise the redirect policy.

### D2. Concurrency with the standard library, results indexed by position

`List*ByTickers` starts one goroutine per ticker; the fetcher's semaphore bounds the real parallelism. Each goroutine writes its result or failure into slot `i` of pre-sized slices, and a `sync.WaitGroup` joins them. The final slices are compacted in index order. That way the output keeps the order the tickers were given, whatever order the responses arrive in.

`PurchaseBalancesBySecurities` runs the stock list and the REIT list in two goroutines and concatenates the results as stocks, then REITs.

- *Alternative:* `golang.org/x/sync/errgroup`. Rejected because it adds a dependency, and errgroup's cancel-on-first-error is the wrong behavior for lists, which must keep going after a failure.
- *Alternative:* a worker pool reading from a channel. Rejected: it needs more code for the same bound.

### D3. Failure types

```go
// internal/resource
type TickerFailure struct {
    Ticker string
    Type   SecurityType // Stock or REIT
    Err    error
}

type QuoteFailuresError struct{ Failures []*TickerFailure }  // implements error
```

- Scraping and service: `List*ByTickers(tickers) ([]*Security, []*TickerFailure)`.
- Service: `PurchaseBalancesBySecurities(...) (*PurchaseBalance, error)`. It returns a `*QuoteFailuresError` when any list has failures and does not call `PurchaseBalance` in that case.
- `cmd/list`: prints the table, then one `Warning: <ticker>: <reason>` line per failure to `cmd.ErrOrStderr()`.
- `cmd/purchase-balance`: returns the error. Cobra already prints `Error: ...` to stderr and the app exits with 1. The message lists every failed ticker with its reason, one per line.

`errors.Join` was considered. Rejected because the `cmd` layer needs the ticker of each failure as data, not only as text.

### D4. Let the transport decompress

Remove the hand-written `Accept-Encoding` header and the gzip branch from the fetcher. Go's transport then requests gzip, decodes it transparently and returns a read error for an invalid gzip body, which keeps the existing "Gzip inválido" scenarios true. `deflate` is no longer advertised. The `Connection: keep-alive` header is also removed: it is the default and has no effect.

### D5. Body limit after decompression

`io.ReadAll(io.LimitReader(body, maxBodySize+1))`, with `maxBodySize = 10 << 20`. If more than `maxBodySize` bytes are read, the fetch is an error. The limit applies to decompressed bytes, which is what stops a gzip bomb.

### D6. Redirect policy

`CheckRedirect` allows a hop only when the target scheme is `https` and the target host equals the host of the first request in the chain (`via[0].URL.Host`), keeping Go's limit of 10 hops. A rejected redirect returns an error, so the fetch fails. It does not return the 3xx page as content.

### D7. Ticker length

`tickerPattern` becomes `^[A-Za-z0-9]{1,12}$`. B3 tickers have 5 to 6 characters, so 12 leaves room for other formats without allowing pathological input.

### D8. Amount validation in `cmd`

A helper `validateAmount(float64) error` in `cmd/common.go`, called first in each `purchase-balance` `RunE`, before any service call. A value is valid only if it is not NaN, not infinite, greater than 0 and no larger than `1e12`. An invalid amount is a usage error: usage is not silenced, the message goes to stderr and the exit code is 1. The ceiling keeps `int(amount / price)` far from overflow for any realistic price.

### D9. Output sanitization in `cmd`

A helper in `internal/tools`:

- `CleanText(s)` removes every rune for which `unicode.IsControl` is true, which includes ESC, BEL and C1 controls. This is enough to break any escape sequence.
- `CSVSafe(s)` prefixes `'` when the cleaned text starts with `=`, `+`, `-` or `@`.

`cmd` builds every row through a small `cell(text string)` function. It always applies `CleanText` to provider text (name, document, admin, segment) and applies `CSVSafe` only when `--csv` is set. Warning lines get `CleanText` too. Numbers and dates are produced by the program, so they are not touched.

A go-pretty column `Transformer` was considered. Rejected: it is configured per column name in each table, so a new column could easily be added without it.

### D10. Stop extracting the stock description

Remove the `//div/p[not(@*)]` query from `GetStockByTicker`. `Security.Description` stays in the struct, empty, so the JSON shape of the type does not change.

### D11. CI and release workflows

- `build_and_release.yaml`:
  - `TAG_NAME` is set in the job's `env:` and used only as `"$TAG_NAME"` inside `run:`. No `${{ ... }}` appears in shell text.
  - Top-level `permissions: contents: read`; `contents: write` only on the build jobs and the checksum job.
  - The `test` job also runs `govulncheck`.
  - A final `checksums` job, which needs every build job, downloads the release assets with `gh release download`, writes `SHA256SUMS` and uploads it.
- `ci.yaml`: `go-version-file: go.mod`.
- Both workflows: `actions/checkout` and `actions/setup-go` pinned by full commit SHA, with the tag in a comment; `govulncheck` pinned to an exact version. The SHAs and the version are resolved when the change is implemented and recorded in the workflows.

## Risks / Trade-offs

- [The provider throttles 4 parallel requests] → 4 is conservative, and the value lives in one constant in `config`. Throttled tickers show up as failures with the HTTP status, so they are no longer silent.
- [Tests depend on timing to prove concurrency] → The limit test counts in-flight requests on an `httptest` server with an atomic counter and a barrier. It never uses sleeps to infer ordering. The ordering test makes the server delay the first ticker longest.
- [The interface changes ripple through the fakes] → The `cmd` and `service` fakes are updated in the same task as the interfaces, so the tree compiles at every step.
- [`purchase-balance` now fails where it used to print a plan] → This is the intended behavior (proposal: BREAKING). The README's error-behavior section is updated with it.
- [Mixed language in the main specs] → The product specs are in Portuguese, while AGENTS.md requires English for OpenSpec artifacts. The deltas keep the existing requirement headers verbatim, as required for matching, and write new text in English, so after archive those specs mix both languages. Translating them is left to a separate change.
- [Pinned SHAs go stale] → That is the cost of pinning. Dependabot for GitHub Actions is the usual answer and can be added later.

## Migration Plan

There is no data or state to migrate. The release notes of the next tag list the three user-visible changes: amount validation, the ticker length limit, and `purchase-balance` aborting on quote failures. Rolling back means reinstalling the previous release.

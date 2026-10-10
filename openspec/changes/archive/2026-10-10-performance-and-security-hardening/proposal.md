# Proposal

## Why

An exploration of the codebase found that almost all of the CLI's runtime is network time spent badly, and that several inputs and outputs are not guarded:

- Quotes are fetched one ticker at a time, stocks before REITs, each with a new HTTP client and a fresh TLS handshake. The wait grows linearly with the number of tickers.
- `--amount inf` makes `purchase-balance` loop forever. Huge or negative amounts produce meaningless counts.
- A ticker whose quote fails silently disappears from a purchase balance, and its share of the money is spent on the other assets. The user acts on an incomplete plan without knowing.
- The response body is read without a size limit, redirects are followed to any host and scheme, and `deflate` is advertised but never decoded (the price then reads as 0 without an error).
- Text scraped from a third-party page reaches the terminal and CSV files unfiltered (control sequences, spreadsheet formulas).
- The release workflow interpolates the tag name directly into shell scripts while holding `contents: write`.

## What Changes

**Input and output**
- **BREAKING** (minor): `--amount` must be a finite number greater than 0 and no larger than a fixed ceiling; other values are a usage error with exit code 1.
- **BREAKING** (minor): tickers longer than a fixed maximum length are rejected as invalid, without a request.
- Control characters are removed from every value written to the terminal or CSV.
- CSV cells that start with `=`, `+`, `-` or `@` are neutralized so spreadsheets do not evaluate them as formulas.

**Fetching quotes**
- One shared HTTP client for the whole run, reusing connections.
- List queries and purchase balances fetch up to 4 tickers at a time, stocks and REITs together, and still return results in the order the tickers were given.
- Response decompression is left to the HTTP transport; the manual `Accept-Encoding` header and gzip handling go away.
- The response body is capped at a fixed size; a larger response is an error.
- Redirects are followed only to `https` URLs on the provider's host.
- The unused stock description is no longer extracted.

**Failures become visible**
- **BREAKING**: `purchase-balance` (all three forms) aborts with exit code 1 and no table when any requested ticker cannot be quoted, naming each failed ticker and its reason.
- `list` commands keep omitting failed tickers, but now write a warning per failed ticker to stderr. Exit code stays 0 while at least one ticker succeeds.

**CI and release**
- Pass the tag name to shell steps as an environment variable, never by template interpolation.
- Grant `contents: write` only to the jobs that upload assets.
- Pin third-party actions by commit SHA and `govulncheck` by version.
- Run `govulncheck` in the release workflow and build CI with the Go version from `go.mod`.
- Publish a `SHA256SUMS` file with the release assets.

## Capabilities

### New Capabilities
<!-- None. -->

### Modified Capabilities
- `stock-quote`: ticker length limit, provider request rules (shared client, transparent decompression, body limit, redirect policy) and concurrent list queries that report failures.
- `reit-quote`: same changes as `stock-quote` for REITs.
- `purchase-balance`: balancing by tickers fails when any ticker cannot be quoted instead of ignoring it.
- `cli`: `--amount` validation, failure reporting for `list` and `purchase-balance`, and sanitized table/CSV output.

## Impact

- **Code:**
  - `internal/scraping/{common,stock,reit}.go`
  - `internal/service/{stock,reit,purchasebalance}.go`
  - `internal/config/config.go`
  - `internal/common/table.go`, `internal/tools/`
  - `cmd/{stock,reit,security,common}.go`
- **Interfaces:** the `List*ByTickers` methods of the scraping and service layers also return the per-ticker failures; `PurchaseBalancesBySecurities` returns an error. Test mocks in `cmd` and `service` change with them.
- **Tests:** new unit tests for amount validation (including `inf` and `NaN`), concurrent ordering, failure reporting, body limit, redirect policy and output sanitization. All offline, using `httptest`.
- **CI/release:** `.github/workflows/ci.yaml`, `.github/workflows/build_and_release.yaml`.
- **Dependencies:** none new; concurrency uses the standard library.
- **User:** invalid amounts and over-long tickers become usage errors; a purchase balance with any failed ticker no longer prints a plan; `list` prints warnings on stderr.

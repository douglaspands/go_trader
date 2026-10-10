package scraping

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
	"trader/internal/resource"
)

func newTestStockScraping(url string) *stockScraping {
	config := newFakeConfig()
	scraping := NewStockScraping(config, NewFetcher(config)).(*stockScraping)
	scraping.url = url
	return scraping
}

func tickersOf(securities []*resource.Security) []string {
	tickers := make([]string, 0, len(securities))
	for _, security := range securities {
		tickers = append(tickers, security.Ticker)
	}
	return tickers
}

func TestGetStockByTickerFound(t *testing.T) {
	// GIVEN
	server, _ := serveFixture(t, "stock_full.html")

	// WHEN
	result, err := newTestStockScraping(server.URL).GetStockByTicker("petr4")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Ticker != "PETR4" {
		t.Errorf(`expected ticker "PETR4" and received "%s"`, result.Ticker)
	}
	if result.Type != resource.STOCK_TYPE {
		t.Errorf(`expected type "%s" and received "%s"`, resource.STOCK_TYPE, result.Type)
	}
	if result.Name != "PETROBRAS" {
		t.Errorf(`expected name "PETROBRAS" and received "%s"`, result.Name)
	}
	if result.Document != "33.000.167/0001-01" {
		t.Errorf(`unexpected document "%s"`, result.Document)
	}
	if result.Description != "" {
		t.Errorf(`expected empty description and received "%s"`, result.Description)
	}
	if result.Price != 1234.56 {
		t.Errorf("expected price 1234.56 and received %f", result.Price)
	}
	if result.Currency.String() != "R$ BRL" || result.Currency.Description != "Brazilian Real" {
		t.Errorf("unexpected currency %+v", result.Currency)
	}
	if result.Origin != server.URL+"/acoes/petr4" {
		t.Errorf(`unexpected origin "%s"`, result.Origin)
	}
	if result.CapturedAt.IsZero() {
		t.Error("expected captured at to be set")
	}
}

func TestGetStockByTickerNameAndPrice(t *testing.T) {
	cases := []struct {
		name         string
		fixture      string
		expectedName string
		expectedCost float64
	}{
		{"hyphen in name", "stock_hyphen_name.html", "FOO-BAR S.A.", 10.5},
		{"title without separator", "stock_no_dash.html", "", 10.5},
		{"price with spaces", "stock_price_spaces.html", "PETROBRAS", 38.15},
		{"price missing", "stock_no_price.html", "PETROBRAS", 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// GIVEN
			server, _ := serveFixture(t, c.fixture)

			// WHEN
			result, err := newTestStockScraping(server.URL).GetStockByTicker("XYZW3")

			// THEN
			if err != nil {
				t.Fatalf("expected no error, received %v", err)
			}
			if result.Name != c.expectedName {
				t.Errorf(`expected name "%s" and received "%s"`, c.expectedName, result.Name)
			}
			if result.Price != c.expectedCost {
				t.Errorf("expected price %f and received %f", c.expectedCost, result.Price)
			}
		})
	}
}

func TestGetStockByTickerEmptyDocument(t *testing.T) {
	// GIVEN
	server, _ := newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body></body></html>"))
	})

	// WHEN
	result, err := newTestStockScraping(server.URL).GetStockByTicker("PETR4")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Name != "" || result.Document != "" || result.Price != 0 {
		t.Errorf("expected empty data and received %+v", result)
	}
}

func TestGetStockByTickerInvalidTicker(t *testing.T) {
	server, requests := newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request to %s", r.URL)
	})
	for _, ticker := range []string{"", "PE%ZZ", "../x", "MX RF11", "ABCDEFGHIJK13"} {
		t.Run(ticker, func(t *testing.T) {
			// WHEN
			result, err := newTestStockScraping(server.URL).GetStockByTicker(ticker)

			// THEN
			if err == nil || !strings.Contains(err.Error(), "invalid ticker") {
				t.Errorf("expected an invalid ticker error, received %v", err)
			}
			if result != nil {
				t.Errorf("expected no stock, received %+v", result)
			}
		})
	}
	if requests.Load() != 0 {
		t.Errorf("expected no requests, received %d", requests.Load())
	}
}

func TestGetStockByTickerMaximumLength(t *testing.T) {
	// GIVEN
	server, requests := serveFixture(t, "stock_full.html")

	// WHEN
	result, err := newTestStockScraping(server.URL).GetStockByTicker("ABCDEFGHIJ12")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Ticker != "ABCDEFGHIJ12" {
		t.Errorf(`expected ticker "ABCDEFGHIJ12" and received "%s"`, result.Ticker)
	}
	if requests.Load() != 1 {
		t.Errorf("expected 1 request, received %d", requests.Load())
	}
}

func TestGetStockByTickerFailures(t *testing.T) {
	cases := []struct {
		name   string
		server func(t *testing.T) string
		expect string
	}{
		{"not found", func(t *testing.T) string {
			server, _ := serveStatus(t, 404)
			return server.URL
		}, `status="404"`},
		{"server error", func(t *testing.T) string {
			server, _ := serveStatus(t, 500)
			return server.URL
		}, `status="500"`},
		{"invalid gzip", func(t *testing.T) string {
			server, _ := serveInvalidGzip(t)
			return server.URL
		}, ""},
		{"unparsable html", func(t *testing.T) string {
			server, _ := serveUnparsable(t)
			return server.URL
		}, ""},
		{"connection refused", func(t *testing.T) string {
			server, _ := serveStatus(t, 200)
			server.Close()
			return server.URL
		}, ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// WHEN
			result, err := newTestStockScraping(c.server(t)).GetStockByTicker("XXXX0")

			// THEN
			if err == nil {
				t.Fatal("expected error, received nil")
			}
			if !strings.Contains(err.Error(), c.expect) {
				t.Errorf(`expected error to contain "%s" and received "%v"`, c.expect, err)
			}
			if c.expect != "" && !strings.Contains(err.Error(), "/acoes/xxxx0") {
				t.Errorf("expected error to mention the url and received %v", err)
			}
			if result != nil {
				t.Errorf("expected no stock, received %+v", result)
			}
		})
	}
}

func TestGetStockByTickerBodyLimit(t *testing.T) {
	t.Run("exactly the limit", func(t *testing.T) {
		server, _ := serveSize(t, maxBodySize)
		result, err := newTestStockScraping(server.URL).GetStockByTicker("PETR4")
		if err != nil || result == nil {
			t.Errorf("expected a stock and no error, received %v and %v", result, err)
		}
	})
	t.Run("one byte over the limit", func(t *testing.T) {
		server, _ := serveSize(t, maxBodySize+1)
		result, err := newTestStockScraping(server.URL).GetStockByTicker("PETR4")
		if err == nil || !strings.Contains(err.Error(), "body larger than") {
			t.Errorf("expected a body size error, received %v", err)
		}
		if result != nil {
			t.Errorf("expected no stock, received %+v", result)
		}
	})
}

func TestGetStockByTickerGzip(t *testing.T) {
	// GIVEN
	server, _ := serveGzip(t, readFixture(t, "stock_full.html"))

	// WHEN
	result, err := newTestStockScraping(server.URL).GetStockByTicker("PETR4")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Name != "PETROBRAS" || result.Price != 1234.56 {
		t.Errorf("unexpected stock %+v", result)
	}
}

func TestListStocksByTickersReportsFailures(t *testing.T) {
	// GIVEN
	server, requests := serveByPath(t, map[string]string{
		"/acoes/petr4": "stock_full.html",
		"/acoes/vale3": "stock_full.html",
	})

	// WHEN
	result, failures := newTestStockScraping(server.URL).ListStocksByTickers([]string{"PETR4", "XXXX0", "BAD TICKER", "VALE3"})

	// THEN
	tickers := tickersOf(result)
	if len(tickers) != 2 || tickers[0] != "PETR4" || tickers[1] != "VALE3" {
		t.Errorf("expected [PETR4 VALE3] and received %v", tickers)
	}
	if len(failures) != 2 {
		t.Fatalf("expected 2 failures and received %d", len(failures))
	}
	if failures[0].Ticker != "XXXX0" || failures[0].Type != resource.STOCK_TYPE || !strings.Contains(failures[0].Err.Error(), `status="404"`) {
		t.Errorf("unexpected failure %+v", failures[0])
	}
	if failures[1].Ticker != "BAD TICKER" || !strings.Contains(failures[1].Err.Error(), "invalid ticker") {
		t.Errorf("unexpected failure %+v", failures[1])
	}
	if requests.Load() != 3 {
		t.Errorf("expected 3 requests and received %d", requests.Load())
	}
}

func TestListStocksByTickersKeepsOrder(t *testing.T) {
	// GIVEN a server that answers C1 first, then B1, and A1 last
	content := readFixture(t, "stock_full.html")
	answered := map[string]chan struct{}{"/acoes/c1": make(chan struct{}), "/acoes/b1": make(chan struct{}), "/acoes/a1": make(chan struct{})}
	waitFor := map[string]string{"/acoes/b1": "/acoes/c1", "/acoes/a1": "/acoes/b1"}
	server, _ := newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		if previous, ok := waitFor[r.URL.Path]; ok {
			select {
			case <-answered[previous]:
			case <-time.After(5 * time.Second):
				t.Errorf("%s waited too long for %s", r.URL.Path, previous)
			}
		}
		w.Write(content)
		w.(http.Flusher).Flush()
		close(answered[r.URL.Path])
	})

	// WHEN
	result, failures := newTestStockScraping(server.URL).ListStocksByTickers([]string{"A1", "B1", "C1"})

	// THEN
	if len(failures) != 0 {
		t.Errorf("expected no failures and received %v", failures)
	}
	if tickers := tickersOf(result); !reflect.DeepEqual(tickers, []string{"A1", "B1", "C1"}) {
		t.Errorf("expected [A1 B1 C1] and received %v", tickers)
	}
}

func TestListStocksByTickersLimitsRequestsInFlight(t *testing.T) {
	// GIVEN
	server, inFlight := serveInFlight(t, readFixture(t, "stock_full.html"))
	tickers := make([]string, 10)
	for i := range tickers {
		tickers[i] = fmt.Sprintf("TICK%d", i)
	}

	// WHEN
	result, failures := newTestStockScraping(server.URL).ListStocksByTickers(tickers)

	// THEN
	if len(result) != 10 || len(failures) != 0 {
		t.Errorf("expected 10 stocks and no failures, received %d and %v", len(result), failures)
	}
	if inFlight.max.Load() != 4 {
		t.Errorf("expected at most and at least 4 requests in flight, received %d", inFlight.max.Load())
	}
}

func TestNewStockScraping(t *testing.T) {
	// WHEN
	fetcher := NewFetcher(newFakeConfig())
	scraping := NewStockScraping(newFakeConfig(), fetcher).(*stockScraping)

	// THEN
	if scraping.url != STATUS_INVEST_URL {
		t.Errorf(`unexpected url "%s"`, scraping.url)
	}
	if scraping.fetcher != fetcher {
		t.Error("expected the given fetcher")
	}
}

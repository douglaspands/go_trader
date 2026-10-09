package scraping

import (
	"net/http"
	"strings"
	"testing"
	"trader/internal/resource"
)

func newTestStockScraping(url string) *stockScraping {
	return &stockScraping{url: url, config: newFakeConfig()}
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
	if result.Description != "Primeiro parágrafo. Segundo parágrafo." {
		t.Errorf(`unexpected description "%s"`, result.Description)
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
	for _, ticker := range []string{"", "PE%ZZ", "../x", "MX RF11"} {
		t.Run(ticker, func(t *testing.T) {
			// WHEN
			result, err := newTestStockScraping(server.URL).GetStockByTicker(ticker)

			// THEN
			if err == nil {
				t.Error("expected error, received nil")
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

func TestListStocksByTickersSkipsFailures(t *testing.T) {
	// GIVEN
	server, requests := serveByPath(t, map[string]string{
		"/acoes/petr4": "stock_full.html",
		"/acoes/vale3": "stock_full.html",
	})

	// WHEN
	result := newTestStockScraping(server.URL).ListStocksByTickers([]string{"PETR4", "XXXX0", "BAD TICKER", "VALE3"})

	// THEN
	tickers := tickersOf(result)
	if len(tickers) != 2 || tickers[0] != "PETR4" || tickers[1] != "VALE3" {
		t.Errorf("expected [PETR4 VALE3] and received %v", tickers)
	}
	if requests.Load() != 3 {
		t.Errorf("expected 3 requests and received %d", requests.Load())
	}
}

func TestNewStockScraping(t *testing.T) {
	// WHEN
	scraping := NewStockScraping(newFakeConfig()).(*stockScraping)

	// THEN
	if scraping.url != STATUS_INVEST_URL {
		t.Errorf(`unexpected url "%s"`, scraping.url)
	}
}

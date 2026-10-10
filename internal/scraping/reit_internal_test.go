package scraping

import (
	"net/http"
	"strings"
	"sync"
	"testing"
	"trader/internal/resource"
)

func newTestReitScraping(url string) *reitScraping {
	config := newFakeConfig()
	scraping := NewReitScraping(config, NewFetcher(config)).(*reitScraping)
	scraping.url = url
	return scraping
}

func TestGetReitByTickerFound(t *testing.T) {
	// GIVEN
	server, _ := serveFixture(t, "reit_full.html")

	// WHEN
	result, err := newTestReitScraping(server.URL).GetReitByTicker("mxrf11")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Ticker != "MXRF11" {
		t.Errorf(`expected ticker "MXRF11" and received "%s"`, result.Ticker)
	}
	if result.Type != resource.REIT_TYPE {
		t.Errorf(`expected type "%s" and received "%s"`, resource.REIT_TYPE, result.Type)
	}
	if result.Name != "Maxi Renda FII" {
		t.Errorf(`unexpected name "%s"`, result.Name)
	}
	if result.Admin != "BANCO XP S.A." {
		t.Errorf(`unexpected admin "%s"`, result.Admin)
	}
	if result.Document != "97.521.225/0001-25" {
		t.Errorf(`unexpected document "%s"`, result.Document)
	}
	if result.Segment != "Híbrido" {
		t.Errorf(`unexpected segment "%s"`, result.Segment)
	}
	if result.Price != 1050.75 {
		t.Errorf("expected price 1050.75 and received %f", result.Price)
	}
	if result.Currency.String() != "R$ BRL" || result.Currency.Description != "Brazilian Real" {
		t.Errorf("unexpected currency %+v", result.Currency)
	}
	if result.Origin != server.URL+"/fundos-imobiliarios/mxrf11" {
		t.Errorf(`unexpected origin "%s"`, result.Origin)
	}
	if result.CapturedAt.IsZero() {
		t.Error("expected captured at to be set")
	}
}

func TestGetReitByTickerWithoutSegment(t *testing.T) {
	// GIVEN
	server, _ := serveFixture(t, "reit_no_segment.html")

	// WHEN
	result, err := newTestReitScraping(server.URL).GetReitByTicker("MXRF11")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Segment != "" {
		t.Errorf(`expected empty segment and received "%s"`, result.Segment)
	}
	if result.Name != "Maxi Renda FII" || result.Admin != "BANCO XP S.A." || result.Document != "97.521.225/0001-25" || result.Price != 1050.75 {
		t.Errorf("unexpected reit %+v", result)
	}
}

func TestGetReitByTickerEmptyDocument(t *testing.T) {
	// GIVEN
	server, _ := newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body></body></html>"))
	})

	// WHEN
	result, err := newTestReitScraping(server.URL).GetReitByTicker("MXRF11")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Name != "" || result.Admin != "" || result.Document != "" || result.Segment != "" || result.Price != 0 {
		t.Errorf("expected empty data and received %+v", result)
	}
}

func TestGetReitByTickerNotFound(t *testing.T) {
	// GIVEN
	server, _ := serveFixture(t, "reit_not_found.html")

	// WHEN
	result, err := newTestReitScraping(server.URL).GetReitByTicker("XXXX00")

	// THEN
	if err == nil {
		t.Fatal("expected error, received nil")
	}
	if !strings.Contains(err.Error(), server.URL+"/fundos-imobiliarios/xxxx00") {
		t.Errorf("expected error to mention the url and received %v", err)
	}
	if result != nil {
		t.Errorf("expected no reit, received %+v", result)
	}
}

func TestGetReitByTickerInvalidTicker(t *testing.T) {
	server, requests := newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request to %s", r.URL)
	})
	for _, ticker := range []string{"", "%zz", "MX RF11", "ABCDEFGHIJK13"} {
		t.Run(ticker, func(t *testing.T) {
			// WHEN
			result, err := newTestReitScraping(server.URL).GetReitByTicker(ticker)

			// THEN
			if err == nil || !strings.Contains(err.Error(), "invalid ticker") {
				t.Errorf("expected an invalid ticker error, received %v", err)
			}
			if result != nil {
				t.Errorf("expected no reit, received %+v", result)
			}
		})
	}
	if requests.Load() != 0 {
		t.Errorf("expected no requests, received %d", requests.Load())
	}
}

func TestGetReitByTickerMaximumLength(t *testing.T) {
	// GIVEN
	server, requests := serveFixture(t, "reit_full.html")

	// WHEN
	result, err := newTestReitScraping(server.URL).GetReitByTicker("ABCDEFGHIJ12")

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

func TestGetReitByTickerFailures(t *testing.T) {
	cases := []struct {
		name   string
		server func(t *testing.T) string
	}{
		{"server error", func(t *testing.T) string {
			server, _ := serveStatus(t, 500)
			return server.URL
		}},
		{"invalid gzip", func(t *testing.T) string {
			server, _ := serveInvalidGzip(t)
			return server.URL
		}},
		{"unparsable html", func(t *testing.T) string {
			server, _ := serveUnparsable(t)
			return server.URL
		}},
		{"connection refused", func(t *testing.T) string {
			server, _ := serveStatus(t, 200)
			server.Close()
			return server.URL
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// WHEN
			result, err := newTestReitScraping(c.server(t)).GetReitByTicker("MXRF11")

			// THEN
			if err == nil {
				t.Fatal("expected error, received nil")
			}
			if result != nil {
				t.Errorf("expected no reit, received %+v", result)
			}
		})
	}
}

func TestGetReitByTickerBodyLimit(t *testing.T) {
	t.Run("exactly the limit", func(t *testing.T) {
		server, _ := serveSize(t, maxBodySize)
		result, err := newTestReitScraping(server.URL).GetReitByTicker("MXRF11")
		if err != nil || result == nil {
			t.Errorf("expected a reit and no error, received %v and %v", result, err)
		}
	})
	t.Run("one byte over the limit", func(t *testing.T) {
		server, _ := serveSize(t, maxBodySize+1)
		result, err := newTestReitScraping(server.URL).GetReitByTicker("MXRF11")
		if err == nil || !strings.Contains(err.Error(), "body larger than") {
			t.Errorf("expected a body size error, received %v", err)
		}
		if result != nil {
			t.Errorf("expected no reit, received %+v", result)
		}
	})
}

func TestGetReitByTickerGzip(t *testing.T) {
	// GIVEN
	server, _ := serveGzip(t, readFixture(t, "reit_full.html"))

	// WHEN
	result, err := newTestReitScraping(server.URL).GetReitByTicker("MXRF11")

	// THEN
	if err != nil {
		t.Fatalf("expected no error, received %v", err)
	}
	if result.Name != "Maxi Renda FII" || result.Price != 1050.75 {
		t.Errorf("unexpected reit %+v", result)
	}
}

func TestListReitsByTickersReportsFailures(t *testing.T) {
	// GIVEN
	server, _ := serveByPath(t, map[string]string{
		"/fundos-imobiliarios/mxrf11": "reit_full.html",
		"/fundos-imobiliarios/xxxx00": "reit_not_found.html",
		"/fundos-imobiliarios/hglg11": "reit_full.html",
	})

	// WHEN
	result, failures := newTestReitScraping(server.URL).ListReitsByTickers([]string{"MXRF11", "XXXX00", "HGLG11"})

	// THEN
	tickers := tickersOf(result)
	if len(tickers) != 2 || tickers[0] != "MXRF11" || tickers[1] != "HGLG11" {
		t.Errorf("expected [MXRF11 HGLG11] and received %v", tickers)
	}
	if len(failures) != 1 {
		t.Fatalf("expected 1 failure and received %d", len(failures))
	}
	if failures[0].Ticker != "XXXX00" || failures[0].Type != resource.REIT_TYPE || !strings.Contains(failures[0].Err.Error(), "ticker not found") {
		t.Errorf("unexpected failure %+v", failures[0])
	}
}

func TestListReitsByTickersLimitsRequestsInFlight(t *testing.T) {
	// GIVEN
	server, inFlight := serveInFlight(t, readFixture(t, "reit_full.html"))
	tickers := []string{"R0", "R1", "R2", "R3", "R4", "R5", "R6", "R7", "R8", "R9"}

	// WHEN
	result, failures := newTestReitScraping(server.URL).ListReitsByTickers(tickers)

	// THEN
	if len(result) != 10 || len(failures) != 0 {
		t.Errorf("expected 10 reits and no failures, received %d and %v", len(result), failures)
	}
	if inFlight.max.Load() != 4 {
		t.Errorf("expected at most and at least 4 requests in flight, received %d", inFlight.max.Load())
	}
}

func TestStocksAndReitsShareTheLimit(t *testing.T) {
	// GIVEN one fetcher shared by both scrapers
	content := append(readFixture(t, "stock_full.html"), readFixture(t, "reit_full.html")...)
	server, inFlight := serveInFlight(t, content)
	config := newFakeConfig()
	fetcher := NewFetcher(config)
	stocks := NewStockScraping(config, fetcher).(*stockScraping)
	stocks.url = server.URL
	reits := NewReitScraping(config, fetcher).(*reitScraping)
	reits.url = server.URL

	// WHEN
	var wg sync.WaitGroup
	wg.Go(func() { stocks.ListStocksByTickers([]string{"S1", "S2", "S3"}) })
	wg.Go(func() { reits.ListReitsByTickers([]string{"R1", "R2", "R3"}) })
	wg.Wait()

	// THEN
	if inFlight.max.Load() != 4 {
		t.Errorf("expected at most and at least 4 requests in flight, received %d", inFlight.max.Load())
	}
}

func TestNewReitScraping(t *testing.T) {
	// WHEN
	fetcher := NewFetcher(newFakeConfig())
	scraping := NewReitScraping(newFakeConfig(), fetcher).(*reitScraping)

	// THEN
	if scraping.url != STATUS_INVEST_URL {
		t.Errorf(`unexpected url "%s"`, scraping.url)
	}
	if scraping.fetcher != fetcher {
		t.Error("expected the given fetcher")
	}
}

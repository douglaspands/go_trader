package scraping

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"sync"
	"trader/internal/config"
	"trader/internal/resource"
)

var tickerPattern = regexp.MustCompile(`^[A-Za-z0-9]{1,12}$`)

func validateTicker(ticker string) error {
	if !tickerPattern.MatchString(ticker) {
		return fmt.Errorf("invalid ticker=\"%s\"", ticker)
	}
	return nil
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/106.0.0.0",
	"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/20.0 Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Android 10; Mobile; rv:120.0) Gecko/120.0 Firefox/120.0",
	"Opera/9.80 (Android; Opera Mini/82.0.2254/123.117; U; en) Presto/2.12.423 Version/12.16",
	"Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/120.0.0.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (iPad; CPU OS 16_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 10; SM-T870) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 EdgA/120.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/18.1.1 Safari/605.1.15",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/132.0.0.0 Safari/537.36 Edg/132.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36 Edg/128.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:134.0) Gecko/20100101 Firefox/134.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 Edg/131.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Windows NT 6.1; rv:109.0) Gecko/20100101 Firefox/115.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36 OPR/116.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36 OPR/95.0.0.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:131.0) Gecko/20100101 Firefox/131.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36",
}

func getUserAgent() string {
	index := rand.Intn(len(userAgents))
	return userAgents[index]
}

// setHeaders sets browser headers. Accept-Encoding is left to the transport, which then
// decompresses the response transparently.
func setHeaders(header *http.Header) {
	header.Set("User-Agent", getUserAgent())
	header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
	header.Set("Cache-Control", "no-cache")
	header.Set("Referer", "https://www.google.com/")
}

// maxBodySize is the largest response body accepted, counted after decompression.
const maxBodySize = 10 << 20

// maxRedirects is the largest number of redirects followed by one request.
const maxRedirects = 10

// checkRedirect follows a redirect only to https on the host of the first request.
func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxRedirects {
		return fmt.Errorf("stopped after %d redirects", maxRedirects)
	}
	if req.URL.Scheme != "https" || req.URL.Host != via[0].URL.Host {
		return fmt.Errorf("redirect to url=\"%s\" not allowed", req.URL)
	}
	return nil
}

// Fetcher downloads pages from the provider. It shares one HTTP client and one limit of
// requests in flight between every scraper built with it.
type Fetcher struct {
	client *http.Client
	slots  chan struct{}
}

type FetcherOption func(*Fetcher)

// WithTransport replaces the HTTP transport. It is meant for tests.
func WithTransport(transport http.RoundTripper) FetcherOption {
	return func(f *Fetcher) {
		f.client.Transport = transport
	}
}

func NewFetcher(config config.Config, options ...FetcherOption) *Fetcher {
	f := &Fetcher{
		client: &http.Client{
			Timeout:       config.GetScrapingTimeout(),
			CheckRedirect: checkRedirect,
		},
		slots: make(chan struct{}, max(config.GetMaxConcurrentRequests(), 1)),
	}
	for _, option := range options {
		option(f)
	}
	return f
}

// listByTickers queries every ticker concurrently, bounded by the fetcher's limit. It returns
// the securities obtained in the order of tickers and one failure per ticker that failed.
func listByTickers(tickers []string, securityType string, get func(string) (*resource.Security, error)) ([]*resource.Security, []*resource.TickerFailure) {
	results := make([]*resource.Security, len(tickers))
	errs := make([]error, len(tickers))
	var wg sync.WaitGroup
	for i, ticker := range tickers {
		wg.Go(func() {
			results[i], errs[i] = get(ticker)
		})
	}
	wg.Wait()

	securities := make([]*resource.Security, 0, len(tickers))
	var failures []*resource.TickerFailure
	for i, ticker := range tickers {
		if errs[i] != nil {
			failures = append(failures, &resource.TickerFailure{Ticker: ticker, Type: securityType, Err: errs[i]})
			continue
		}
		securities = append(securities, results[i])
	}
	return securities, failures
}

// Fetch returns the body of url. A slot is held only while the request and the body read
// are in progress, so parsing does not count against the limit.
func (f *Fetcher) Fetch(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(&req.Header)

	f.slots <- struct{}{}
	defer func() { <-f.slots }()

	httpResponse, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()
	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("status=\"%d\" for url=\"%s\"", httpResponse.StatusCode, url)
	}
	body, err := io.ReadAll(io.LimitReader(httpResponse.Body, maxBodySize+1))
	if err != nil {
		return nil, err
	}
	if len(body) > maxBodySize {
		return nil, fmt.Errorf("body larger than %d bytes for url=\"%s\"", maxBodySize, url)
	}
	return body, nil
}

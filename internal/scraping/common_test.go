package scraping

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetUserAgent(t *testing.T) {
	ua := getUserAgent()
	if ua == "" {
		t.Error("Expected user agent to be non-empty")
	}
}

// The provider answers 403 to a User-Agent whose last version number is truncated.
func TestUserAgentsAreComplete(t *testing.T) {
	truncated := regexp.MustCompile(`(\.|Safari/537\.3|Safari/605\.1\.1)$`)
	for _, ua := range userAgents {
		if truncated.MatchString(ua) {
			t.Errorf("truncated user agent %q", ua)
		}
	}
}

func TestSetHeaders(t *testing.T) {
	header := make(http.Header)
	setHeaders(&header)

	if header.Get("User-Agent") == "" {
		t.Error("Expected User-Agent header to be set")
	}
	if header.Get("Accept") == "" {
		t.Error("Expected Accept header to be set")
	}
	if header.Get("Accept-Encoding") != "" {
		t.Errorf("Expected Accept-Encoding to be left to the transport, got %s", header.Get("Accept-Encoding"))
	}
	if header.Get("Connection") != "" {
		t.Errorf("Expected no Connection header, got %s", header.Get("Connection"))
	}
}

func newTestFetcher(timeout time.Duration) *Fetcher {
	return NewFetcher(fakeConfig{timeout: timeout})
}

func TestFetch(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body>Hello</body></html>"))
		}))
		defer server.Close()

		content, err := newTestFetcher(time.Second).Fetch(server.URL)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(content) != "<html><body>Hello</body></html>" {
			t.Errorf("Expected content to match, got %s", content)
		}
	})

	t.Run("AdvertisesOnlyGzip", func(t *testing.T) {
		encodings := make(chan string, 1)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encodings <- r.Header.Get("Accept-Encoding")
		}))
		defer server.Close()

		if _, err := newTestFetcher(time.Second).Fetch(server.URL); err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if encoding := <-encodings; encoding != "gzip" {
			t.Errorf("Expected Accept-Encoding gzip, got %q", encoding)
		}
	})

	t.Run("SuccessGzip", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusOK)
			gw := gzip.NewWriter(w)
			gw.Write([]byte("<html><body>Gzip</body></html>"))
			gw.Close()
		}))
		defer server.Close()

		content, err := newTestFetcher(time.Second).Fetch(server.URL)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(content) != "<html><body>Gzip</body></html>" {
			t.Errorf("Expected content to match, got %s", content)
		}
	})

	t.Run("StatusError", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		_, err := newTestFetcher(time.Second).Fetch(server.URL)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if !strings.Contains(err.Error(), "status=\"404\"") {
			t.Errorf("Expected error to contain status 404, got %v", err)
		}
	})

	t.Run("Timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(2 * time.Second)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		// Set timeout shorter than the sleep
		_, err := newTestFetcher(time.Second).Fetch(server.URL)
		if err == nil {
			t.Fatal("Expected timeout error, got nil")
		}
	})

	t.Run("InvalidUrl", func(t *testing.T) {
		content, err := newTestFetcher(time.Second).Fetch("http://example.com/%zz")
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if content != nil {
			t.Errorf("Expected no content, got %s", content)
		}
	})

	t.Run("RequestFailure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		url := server.URL
		server.Close()

		_, err := newTestFetcher(time.Second).Fetch(url)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
	})

	t.Run("InvalidGzip", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("this is not gzip"))
		}))
		defer server.Close()

		content, err := newTestFetcher(time.Second).Fetch(server.URL)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		if content != nil {
			t.Errorf("Expected no content, got %s", content)
		}
	})

	t.Run("BodyClosedOnStatusError", func(t *testing.T) {
		closed := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("partial body"))
			w.(http.Flusher).Flush()
			select {
			case <-r.Context().Done():
				close(closed)
			case <-time.After(5 * time.Second):
			}
		}))
		defer server.Close()

		_, err := newTestFetcher(10 * time.Second).Fetch(server.URL)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}
		select {
		case <-closed:
		case <-time.After(2 * time.Second):
			t.Error("Expected connection to be closed after a status error")
		}
	})
}

func TestFetchBodyLimit(t *testing.T) {
	cases := []struct {
		name    string
		size    int
		gzipped bool
		fails   bool
	}{
		{"exactly the limit", maxBodySize, false, false},
		{"one byte over the limit", maxBodySize + 1, false, true},
		{"gzip over the limit after decompression", maxBodySize + 1, true, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// GIVEN
			body := []byte(strings.Repeat("a", c.size))
			if c.gzipped {
				body = gzipBytes(t, body)
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if c.gzipped {
					w.Header().Set("Content-Encoding", "gzip")
				}
				w.Write(body)
			}))
			defer server.Close()

			// WHEN
			content, err := newTestFetcher(5 * time.Second).Fetch(server.URL)

			// THEN
			if c.fails {
				if err == nil || !strings.Contains(err.Error(), "body larger than") {
					t.Errorf("Expected a body size error, got %v", err)
				}
				if content != nil {
					t.Errorf("Expected no content, got %d bytes", len(content))
				}
				return
			}
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if len(content) != c.size {
				t.Errorf("Expected %d bytes, got %d", c.size, len(content))
			}
		})
	}
}

// newTLSFetcher builds a fetcher that trusts the certificate of server.
func newTLSFetcher(server *httptest.Server) *Fetcher {
	return NewFetcher(newFakeConfig(), WithTransport(server.Client().Transport))
}

func TestFetchRedirect(t *testing.T) {
	t.Run("SameHostHttps", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/from" {
				http.Redirect(w, r, "/to", http.StatusFound)
				return
			}
			w.Write([]byte("arrived at " + r.URL.Path))
		}))
		defer server.Close()

		content, err := newTLSFetcher(server).Fetch(server.URL + "/from")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(content) != "arrived at /to" {
			t.Errorf("Expected the redirect to be followed, got %s", content)
		}
	})

	t.Run("AnotherHost", func(t *testing.T) {
		var otherRequests atomic.Int32
		other := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			otherRequests.Add(1)
		}))
		defer other.Close()
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, other.URL+"/to", http.StatusFound)
		}))
		defer server.Close()

		content, err := newTLSFetcher(server).Fetch(server.URL + "/from")
		if err == nil || !strings.Contains(err.Error(), "not allowed") {
			t.Errorf("Expected a redirect error, got %v", err)
		}
		if content != nil {
			t.Errorf("Expected no content, got %s", content)
		}
		if otherRequests.Load() != 0 {
			t.Errorf("Expected no request to the other host, got %d", otherRequests.Load())
		}
	})

	t.Run("Http", func(t *testing.T) {
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://"+r.Host+"/to", http.StatusFound)
		}))
		defer server.Close()

		content, err := newTLSFetcher(server).Fetch(server.URL + "/from")
		if err == nil || !strings.Contains(err.Error(), "not allowed") {
			t.Errorf("Expected a redirect error, got %v", err)
		}
		if content != nil {
			t.Errorf("Expected no content, got %s", content)
		}
	})

	t.Run("TooManyHops", func(t *testing.T) {
		var requests atomic.Int32
		server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requests.Add(1)
			http.Redirect(w, r, "/again", http.StatusFound)
		}))
		defer server.Close()

		_, err := newTLSFetcher(server).Fetch(server.URL + "/from")
		if err == nil || !strings.Contains(err.Error(), "stopped after 10 redirects") {
			t.Errorf("Expected a redirect limit error, got %v", err)
		}
		if requests.Load() != maxRedirects {
			t.Errorf("Expected %d requests, got %d", maxRedirects, requests.Load())
		}
	})
}

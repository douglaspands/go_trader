package scraping

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type fakeConfig struct {
	timeout time.Duration
}

func (c fakeConfig) GetVersion() string { return "test" }

func (c fakeConfig) GetScrapingTimeout() time.Duration { return c.timeout }

func newFakeConfig() fakeConfig {
	return fakeConfig{timeout: time.Second}
}

func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("reading fixture %s: %v", name, err)
	}
	return content
}

func gzipBytes(t *testing.T, content []byte) []byte {
	t.Helper()
	var buffer bytes.Buffer
	writer := gzip.NewWriter(&buffer)
	if _, err := writer.Write(content); err != nil {
		t.Fatalf("gzip write: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	return buffer.Bytes()
}

// newCountingServer starts a server that counts the requests it receives.
func newCountingServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		handler(w, r)
	}))
	t.Cleanup(server.Close)
	return server, &requests
}

func serveFixture(t *testing.T, name string) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	content := readFixture(t, name)
	return newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	})
}

func serveStatus(t *testing.T, status int) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	return newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
}

func serveGzip(t *testing.T, content []byte) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	compressed := gzipBytes(t, content)
	return newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Write(compressed)
	})
}

func serveInvalidGzip(t *testing.T) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	return newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.Write([]byte("this is not gzip"))
	})
}

// serveUnparsable serves a document nested deeply enough to make the HTML parser fail.
func serveUnparsable(t *testing.T) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	content := strings.Repeat("<div>", 600)
	return newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(content))
	})
}

// serveByPath serves each fixture on its URL path and answers 404 for unknown paths.
func serveByPath(t *testing.T, fixtures map[string]string) (*httptest.Server, *atomic.Int32) {
	t.Helper()
	contents := make(map[string][]byte)
	for path, name := range fixtures {
		contents[path] = readFixture(t, name)
	}
	return newCountingServer(t, func(w http.ResponseWriter, r *http.Request) {
		content, ok := contents[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(content)
	})
}

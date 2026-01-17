package scraping

import (
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetUserAgent(t *testing.T) {
	ua := getUserAgent()
	if ua == "" {
		t.Error("Expected user agent to be non-empty")
	}
}

func TestSetHeaders(t *testing.T) {
	header := make(http.Header)
	setHeaders(&header)

	if header.Get("User-Agent") == "" {
		t.Error("Expected User-Agent header to be set")
	}
	if header.Get("Connection") != "keep-alive" {
		t.Errorf("Expected Connection header to be keep-alive, got %s", header.Get("Connection"))
	}
	if header.Get("Accept") == "" {
		t.Error("Expected Accept header to be set")
	}
}

func TestGetHtml(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("<html><body>Hello</body></html>"))
		}))
		defer server.Close()

		content, err := getHtml(server.URL, 1)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if string(content) != "<html><body>Hello</body></html>" {
			t.Errorf("Expected content to match, got %s", content)
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

		content, err := getHtml(server.URL, 1)
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

		_, err := getHtml(server.URL, 1)
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
		_, err := getHtml(server.URL, 1)
		if err == nil {
			t.Fatal("Expected timeout error, got nil")
		}
	})
}

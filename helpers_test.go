package jupiter

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	return server
}

func newTestClient(serverURL string) *Client {
	client := NewClient(serverURL, "test-api-key")
	client.Limiter = rate.NewLimiter(rate.Inf, 1)
	return client
}

func jsonHandler(t *testing.T, wantMethod, wantPath string, response any) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != wantMethod {
			t.Errorf("expected method %s, got %s", wantMethod, r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("expected path %s, got %s", wantPath, r.URL.Path)
		}
		if r.Header.Get(ApiKeyHeader) != "test-api-key" {
			t.Errorf("expected api key header test-api-key, got %s", r.Header.Get(ApiKeyHeader))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func errorHandler(statusCode int, body string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}
}

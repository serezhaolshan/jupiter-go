package jupiter

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://api.jup.ag", "my-key")

	if client.ApiUrl != "https://api.jup.ag" {
		t.Errorf("expected ApiUrl https://api.jup.ag, got %s", client.ApiUrl)
	}
	if client.ApiKey != "my-key" {
		t.Errorf("expected ApiKey my-key, got %s", client.ApiKey)
	}
	if client.Limiter == nil {
		t.Fatal("expected Limiter to be set")
	}
	if client.c == nil {
		t.Fatal("expected http client to be set")
	}
}

func TestNewClient_DefaultValues(t *testing.T) {
	client := NewClient(DefaultURL, "")

	if client.ApiUrl != DefaultURL {
		t.Errorf("expected ApiUrl %s, got %s", DefaultURL, client.ApiUrl)
	}
	if client.ApiKey != "" {
		t.Errorf("expected empty ApiKey, got %s", client.ApiKey)
	}
}

func TestClientUrl(t *testing.T) {
	client := NewClient("https://api.jup.ag", "")

	tests := []struct {
		endpoint string
		want     string
	}{
		{"/price/v3", "https://api.jup.ag/price/v3"},
		{"/tokens/v2/trending", "https://api.jup.ag/tokens/v2/trending"},
		{"", "https://api.jup.ag"},
	}

	for _, tt := range tests {
		got := client.Url(tt.endpoint)
		if got != tt.want {
			t.Errorf("Url(%q) = %q, want %q", tt.endpoint, got, tt.want)
		}
	}
}

func TestNewRequest(t *testing.T) {
	params := url.Values{}
	params.Set("key", "value")

	req := NewRequest("/test", params)

	if req.Endpoint != "/test" {
		t.Errorf("expected endpoint /test, got %s", req.Endpoint)
	}
	if req.Method != http.MethodGet {
		t.Errorf("expected method GET, got %s", req.Method)
	}
	if req.QueryParams.Get("key") != "value" {
		t.Errorf("expected query param key=value, got %s", req.QueryParams.Get("key"))
	}
}

func TestNewRequest_WithMethod(t *testing.T) {
	req := NewRequest("/test", nil, http.MethodPost)

	if req.Method != http.MethodPost {
		t.Errorf("expected method POST, got %s", req.Method)
	}
}

func TestNewPostRequest(t *testing.T) {
	body := ExecuteRequest{
		SignedTransaction: "tx123",
		RequestID:         "req456",
	}

	req, err := NewPostRequest("/test", body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Endpoint != "/test" {
		t.Errorf("expected endpoint /test, got %s", req.Endpoint)
	}
	if req.Method != http.MethodPost {
		t.Errorf("expected method POST, got %s", req.Method)
	}
	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	bodyBytes, _ := io.ReadAll(req.Body)
	var decoded ExecuteRequest
	json.Unmarshal(bodyBytes, &decoded)
	if decoded.SignedTransaction != "tx123" {
		t.Errorf("expected SignedTransaction tx123, got %s", decoded.SignedTransaction)
	}
}

func TestNewPostRequest_MarshalError(t *testing.T) {
	// channels cannot be marshaled to JSON
	_, err := NewPostRequest("/test", make(chan int))
	if err == nil {
		t.Fatal("expected error for unmarshalable body")
	}
	if !strings.Contains(err.Error(), "failed to marshal") {
		t.Errorf("expected marshal error, got: %v", err)
	}
}

func TestNewHttpRequest(t *testing.T) {
	req := &Request{
		Endpoint: "https://api.jup.ag/test",
		Method:   http.MethodGet,
	}

	httpReq, err := req.NewHttpRequest(context.Background(), "my-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Method != http.MethodGet {
		t.Errorf("expected GET, got %s", httpReq.Method)
	}
	if httpReq.URL.String() != "https://api.jup.ag/test" {
		t.Errorf("expected URL https://api.jup.ag/test, got %s", httpReq.URL.String())
	}
	if httpReq.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", httpReq.Header.Get("Content-Type"))
	}
	if httpReq.Header.Get(ApiKeyHeader) != "my-key" {
		t.Errorf("expected api key my-key, got %s", httpReq.Header.Get(ApiKeyHeader))
	}
}

func TestNewHttpRequest_NoApiKey(t *testing.T) {
	req := &Request{
		Endpoint: "https://api.jup.ag/test",
		Method:   http.MethodGet,
	}

	httpReq, err := req.NewHttpRequest(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if httpReq.Header.Get(ApiKeyHeader) != "" {
		t.Errorf("expected no api key header, got %s", httpReq.Header.Get(ApiKeyHeader))
	}
}

func TestNewHttpRequest_WithQueryParams(t *testing.T) {
	params := url.Values{}
	params.Set("ids", "SOL,USDC")

	req := &Request{
		Endpoint:    "https://api.jup.ag/price/v3",
		Method:      http.MethodGet,
		QueryParams: params,
	}

	httpReq, err := req.NewHttpRequest(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(httpReq.URL.String(), "ids=SOL") {
		t.Errorf("expected URL to contain query params, got %s", httpReq.URL.String())
	}
}

func TestNewHttpRequest_WithBody(t *testing.T) {
	body := bytes.NewReader([]byte(`{"key":"value"}`))
	req := &Request{
		Endpoint: "https://api.jup.ag/test",
		Method:   http.MethodPost,
		Body:     body,
	}

	httpReq, err := req.NewHttpRequest(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	bodyBytes, _ := io.ReadAll(httpReq.Body)
	if string(bodyBytes) != `{"key":"value"}` {
		t.Errorf("expected body {\"key\":\"value\"}, got %s", string(bodyBytes))
	}
}

func TestDoCall_Success(t *testing.T) {
	server := newTestServer(t, jsonHandler(t, http.MethodGet, "/test", map[string]string{"hello": "world"}))
	client := newTestClient(server.URL)

	req := NewRequest(client.Url("/test"), nil)
	var response map[string]string
	_, err := client.doCall(context.Background(), req, &response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response["hello"] != "world" {
		t.Errorf("expected hello=world, got %v", response)
	}
}

func TestDoCall_ErrorStatus(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, `{"error":"bad request"}`))
	client := newTestClient(server.URL)

	req := NewRequest(client.Url("/test"), nil)
	var response map[string]string
	_, err := client.doCall(context.Background(), req, &response)
	if err == nil {
		t.Fatal("expected error for 400 status")
	}
	if !strings.Contains(err.Error(), "400") {
		t.Errorf("expected error to contain status code 400, got: %v", err)
	}
}

func TestDoCall_ServerError(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusInternalServerError, "internal error"))
	client := newTestClient(server.URL)

	req := NewRequest(client.Url("/test"), nil)
	var response map[string]string
	_, err := client.doCall(context.Background(), req, &response)
	if err == nil {
		t.Fatal("expected error for 500 status")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to contain status code 500, got: %v", err)
	}
}

func TestDoCall_NonPointerResponse(t *testing.T) {
	server := newTestServer(t, jsonHandler(t, http.MethodGet, "/test", "ok"))
	client := newTestClient(server.URL)

	req := NewRequest(client.Url("/test"), nil)
	var response string
	_, err := client.doCall(context.Background(), req, response) // not a pointer
	if err == nil {
		t.Fatal("expected error for non-pointer response")
	}
	if !strings.Contains(err.Error(), "not a pointer") {
		t.Errorf("expected pointer error, got: %v", err)
	}
}

func TestDoCall_InvalidJSON(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	})
	client := newTestClient(server.URL)

	req := NewRequest(client.Url("/test"), nil)
	var response map[string]string
	_, err := client.doCall(context.Background(), req, &response)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "could not decode body") {
		t.Errorf("expected decode error, got: %v", err)
	}
}

func TestDoCall_CancelledContext(t *testing.T) {
	server := newTestServer(t, jsonHandler(t, http.MethodGet, "/test", "ok"))
	client := newTestClient(server.URL)
	// Use a slow rate limiter so context cancellation happens during Wait
	client.Limiter = rate.NewLimiter(rate.Every(10*time.Second), 0)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	req := NewRequest(client.Url("/test"), nil)
	var response string
	_, err := client.doCall(ctx, req, &response)
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}

func TestConstants(t *testing.T) {
	if DefaultURL != "https://api.jup.ag" {
		t.Errorf("expected DefaultURL https://api.jup.ag, got %s", DefaultURL)
	}
	if ApiKeyHeader != "x-api-key" {
		t.Errorf("expected ApiKeyHeader x-api-key, got %s", ApiKeyHeader)
	}
	if RateLimitMilliseconds != 120 {
		t.Errorf("expected RateLimitMilliseconds 120, got %d", RateLimitMilliseconds)
	}
}

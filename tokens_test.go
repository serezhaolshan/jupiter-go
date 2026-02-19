package jupiter

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetTokens(t *testing.T) {
	tokens := []TokenV2{
		{ID: "SOL111", Name: "Solana", Symbol: "SOL", Decimals: 9},
		{ID: "USDC111", Name: "USD Coin", Symbol: "USDC", Decimals: 6},
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tokens/v2/trending" {
			t.Errorf("expected path /tokens/v2/trending, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	})
	client := newTestClient(server.URL)

	result, err := client.GetTokens(context.Background(), GetTokensParams{
		SortBy: "trending",
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 tokens, got %d", len(result))
	}
	if result[0].Symbol != "SOL" {
		t.Errorf("expected first token SOL, got %s", result[0].Symbol)
	}
	if result[1].Symbol != "USDC" {
		t.Errorf("expected second token USDC, got %s", result[1].Symbol)
	}
}

func TestGetTokens_WithInterval(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tokens/v2/volume/1h" {
			t.Errorf("expected path /tokens/v2/volume/1h, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]TokenV2{})
	})
	client := newTestClient(server.URL)

	_, err := client.GetTokens(context.Background(), GetTokensParams{
		SortBy:   "volume",
		Interval: "1h",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTokens_NoLimit(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "" {
			t.Errorf("expected no limit param, got %s", r.URL.Query().Get("limit"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]TokenV2{})
	})
	client := newTestClient(server.URL)

	_, err := client.GetTokens(context.Background(), GetTokensParams{
		SortBy: "trending",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTokens_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusInternalServerError, "error"))
	client := newTestClient(server.URL)

	_, err := client.GetTokens(context.Background(), GetTokensParams{SortBy: "trending"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSearchTokens(t *testing.T) {
	tokens := []TokenV2{
		{ID: "SOL111", Name: "Solana", Symbol: "SOL", Decimals: 9},
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/tokens/v2/search" {
			t.Errorf("expected path /tokens/v2/search, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "SOL" {
			t.Errorf("expected query=SOL, got %s", r.URL.Query().Get("query"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokens)
	})
	client := newTestClient(server.URL)

	result, err := client.SearchTokens(context.Background(), SearchTokensParams{Query: "SOL"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 token, got %d", len(result))
	}
	if result[0].Name != "Solana" {
		t.Errorf("expected token name Solana, got %s", result[0].Name)
	}
}

func TestSearchTokens_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "bad request"))
	client := newTestClient(server.URL)

	_, err := client.SearchTokens(context.Background(), SearchTokensParams{Query: ""})
	if err == nil {
		t.Fatal("expected error")
	}
}

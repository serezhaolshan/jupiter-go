package jupiter

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetTriggerOrders(t *testing.T) {
	resp := GetTriggerOrdersResponse{
		Orders: []TriggerOrder{
			{ID: "order1", User: "user1", InputMint: "SOL", OutputMint: "USDC", Status: "open"},
			{ID: "order2", User: "user1", InputMint: "SOL", OutputMint: "USDC", Status: "filled"},
		},
		HasMoreData: true,
		Page:        1,
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/trigger/v1/getTriggerOrders" {
			t.Errorf("expected path /trigger/v1/getTriggerOrders, got %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("user") != "user1" {
			t.Errorf("expected user=user1, got %s", q.Get("user"))
		}
		if q.Get("orderStatus") != "open" {
			t.Errorf("expected orderStatus=open, got %s", q.Get("orderStatus"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	client := newTestClient(server.URL)

	result, err := client.GetTriggerOrders(context.Background(), GetTriggerOrdersParams{
		User:        "user1",
		OrderStatus: "open",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Orders) != 2 {
		t.Fatalf("expected 2 orders, got %d", len(result.Orders))
	}
	if !result.HasMoreData {
		t.Error("expected HasMoreData to be true")
	}
	if result.Page != 1 {
		t.Errorf("expected page 1, got %d", result.Page)
	}
}

func TestGetTriggerOrders_WithOptionalParams(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("inputMint") != "SOL" {
			t.Errorf("expected inputMint=SOL, got %s", q.Get("inputMint"))
		}
		if q.Get("outputMint") != "USDC" {
			t.Errorf("expected outputMint=USDC, got %s", q.Get("outputMint"))
		}
		if q.Get("page") != "2" {
			t.Errorf("expected page=2, got %s", q.Get("page"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetTriggerOrdersResponse{})
	})
	client := newTestClient(server.URL)

	_, err := client.GetTriggerOrders(context.Background(), GetTriggerOrdersParams{
		User:        "user1",
		OrderStatus: "open",
		InputMint:   "SOL",
		OutputMint:  "USDC",
		Page:        2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTriggerOrders_NoOptionalParams(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("inputMint") != "" {
			t.Errorf("expected no inputMint param, got %s", q.Get("inputMint"))
		}
		if q.Get("outputMint") != "" {
			t.Errorf("expected no outputMint param, got %s", q.Get("outputMint"))
		}
		if q.Get("page") != "" {
			t.Errorf("expected no page param, got %s", q.Get("page"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(GetTriggerOrdersResponse{})
	})
	client := newTestClient(server.URL)

	_, err := client.GetTriggerOrders(context.Background(), GetTriggerOrdersParams{
		User:        "user1",
		OrderStatus: "open",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetTriggerOrders_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "missing user"))
	client := newTestClient(server.URL)

	_, err := client.GetTriggerOrders(context.Background(), GetTriggerOrdersParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

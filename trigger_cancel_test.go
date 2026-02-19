package jupiter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCancelOrder(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trigger/v1/cancelOrder" {
			t.Errorf("expected path /trigger/v1/cancelOrder, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req CancelOrderRequest
		json.Unmarshal(body, &req)
		if req.Maker != "maker1" {
			t.Errorf("expected Maker maker1, got %s", req.Maker)
		}
		if req.Order != "order123" {
			t.Errorf("expected Order order123, got %s", req.Order)
		}
		if req.ComputeUnitPrice != "5000" {
			t.Errorf("expected ComputeUnitPrice 5000, got %s", req.ComputeUnitPrice)
		}

		resp := CancelOrderResponse{
			Transaction: "cancel-tx-789",
			RequestID:   "cancel-req-012",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	client := newTestClient(server.URL)

	result, err := client.CancelOrder(context.Background(), CancelOrderRequest{
		Maker:            "maker1",
		Order:            "order123",
		ComputeUnitPrice: "5000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Transaction != "cancel-tx-789" {
		t.Errorf("expected transaction cancel-tx-789, got %s", result.Transaction)
	}
	if result.RequestID != "cancel-req-012" {
		t.Errorf("expected requestId cancel-req-012, got %s", result.RequestID)
	}
}

func TestCancelOrder_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "order not found"))
	client := newTestClient(server.URL)

	_, err := client.CancelOrder(context.Background(), CancelOrderRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

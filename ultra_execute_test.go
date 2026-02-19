package jupiter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestExecuteUltra(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/ultra/v1/execute" {
			t.Errorf("expected path /ultra/v1/execute, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req ExecuteRequest
		json.Unmarshal(body, &req)
		if req.SignedTransaction != "signed-tx-123" {
			t.Errorf("expected SignedTransaction signed-tx-123, got %s", req.SignedTransaction)
		}
		if req.RequestID != "req-456" {
			t.Errorf("expected RequestID req-456, got %s", req.RequestID)
		}

		resp := ExecuteResponse{
			Status:    "Success",
			Signature: "sig789",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	client := newTestClient(server.URL)

	result, err := client.ExecuteUltra(context.Background(), ExecuteRequest{
		SignedTransaction: "signed-tx-123",
		RequestID:         "req-456",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "Success" {
		t.Errorf("expected status Success, got %s", result.Status)
	}
	if result.Signature != "sig789" {
		t.Errorf("expected signature sig789, got %s", result.Signature)
	}
}

func TestExecuteUltra_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "invalid tx"))
	client := newTestClient(server.URL)

	_, err := client.ExecuteUltra(context.Background(), ExecuteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

package jupiter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestExecuteTrigger(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trigger/v1/execute" {
			t.Errorf("expected path /trigger/v1/execute, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req ExecuteRequest
		json.Unmarshal(body, &req)
		if req.SignedTransaction != "trigger-tx-123" {
			t.Errorf("expected SignedTransaction trigger-tx-123, got %s", req.SignedTransaction)
		}

		resp := ExecuteResponse{
			Status:    "Success",
			Signature: "trigger-sig-456",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	client := newTestClient(server.URL)

	result, err := client.ExecuteTrigger(context.Background(), ExecuteRequest{
		SignedTransaction: "trigger-tx-123",
		RequestID:         "req-789",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "Success" {
		t.Errorf("expected status Success, got %s", result.Status)
	}
	if result.Signature != "trigger-sig-456" {
		t.Errorf("expected signature trigger-sig-456, got %s", result.Signature)
	}
}

func TestExecuteTrigger_WithError(t *testing.T) {
	code := 1001

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		resp := ExecuteResponse{
			Status: "Failed",
			Error:  "insufficient funds",
			Code:   &code,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	client := newTestClient(server.URL)

	result, err := client.ExecuteTrigger(context.Background(), ExecuteRequest{
		SignedTransaction: "bad-tx",
		RequestID:         "req-000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Status != "Failed" {
		t.Errorf("expected status Failed, got %s", result.Status)
	}
	if result.Error != "insufficient funds" {
		t.Errorf("expected error insufficient funds, got %s", result.Error)
	}
	if *result.Code != 1001 {
		t.Errorf("expected code 1001, got %d", *result.Code)
	}
}

func TestExecuteTrigger_ServerError(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusInternalServerError, "server error"))
	client := newTestClient(server.URL)

	_, err := client.ExecuteTrigger(context.Background(), ExecuteRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

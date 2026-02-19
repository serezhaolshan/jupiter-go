package jupiter

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetProgramIDToLabel(t *testing.T) {
	labels := ProgramIDToLabelResponse{
		"675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8": "Raydium",
		"whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc":  "Orca",
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/swap/v1/program-id-to-label" {
			t.Errorf("expected path /swap/v1/program-id-to-label, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(labels)
	})
	client := newTestClient(server.URL)

	result, err := client.GetProgramIDToLabel(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(result))
	}
	if result["675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"] != "Raydium" {
		t.Errorf("expected Raydium label")
	}
	if result["whirLbMiicVdio4qvUfM5KAg6Ct8VwpYzGff3uctyCc"] != "Orca" {
		t.Errorf("expected Orca label")
	}
}

func TestGetProgramIDToLabel_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusInternalServerError, "error"))
	client := newTestClient(server.URL)

	_, err := client.GetProgramIDToLabel(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

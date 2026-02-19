package jupiter

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetRouters(t *testing.T) {
	routers := RoutersResponse{
		{ID: "router1", Name: "Jupiter", Icon: "https://jup.ag/icon.png"},
		{ID: "router2", Name: "Raydium"},
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/ultra/v1/order/routers" {
			t.Errorf("expected path /ultra/v1/order/routers, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(routers)
	})
	client := newTestClient(server.URL)

	result, err := client.GetRouters(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 routers, got %d", len(result))
	}
	if result[0].Name != "Jupiter" {
		t.Errorf("expected first router Jupiter, got %s", result[0].Name)
	}
	if result[1].Icon != "" {
		t.Errorf("expected empty icon for second router, got %s", result[1].Icon)
	}
}

func TestGetRouters_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusInternalServerError, "error"))
	client := newTestClient(server.URL)

	_, err := client.GetRouters(context.Background())
	if err == nil {
		t.Fatal("expected error")
	}
}

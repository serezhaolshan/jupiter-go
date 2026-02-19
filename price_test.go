package jupiter

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPrices(t *testing.T) {
	blockID := int64(123456)
	decimals := 9
	change := 5.5

	prices := PriceV3Response{
		"SOL": {
			USDPrice:       150.25,
			BlockID:        &blockID,
			Decimals:       &decimals,
			PriceChange24h: &change,
		},
		"USDC": {
			USDPrice: 1.0,
		},
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/price/v3" {
			t.Errorf("expected path /price/v3, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("ids") != "SOL,USDC" {
			t.Errorf("expected ids=SOL,USDC, got %s", r.URL.Query().Get("ids"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(prices)
	})
	client := newTestClient(server.URL)

	result, err := client.GetPrices(context.Background(), "SOL,USDC")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	solPrice, ok := result["SOL"]
	if !ok {
		t.Fatal("expected SOL price entry")
	}
	if solPrice.USDPrice != 150.25 {
		t.Errorf("expected SOL price 150.25, got %f", solPrice.USDPrice)
	}
	if *solPrice.BlockID != 123456 {
		t.Errorf("expected block ID 123456, got %d", *solPrice.BlockID)
	}
	if *solPrice.PriceChange24h != 5.5 {
		t.Errorf("expected price change 5.5, got %f", *solPrice.PriceChange24h)
	}

	usdcPrice, ok := result["USDC"]
	if !ok {
		t.Fatal("expected USDC price entry")
	}
	if usdcPrice.USDPrice != 1.0 {
		t.Errorf("expected USDC price 1.0, got %f", usdcPrice.USDPrice)
	}
}

func TestGetPrices_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "invalid ids"))
	client := newTestClient(server.URL)

	_, err := client.GetPrices(context.Background(), "INVALID")
	if err == nil {
		t.Fatal("expected error")
	}
}

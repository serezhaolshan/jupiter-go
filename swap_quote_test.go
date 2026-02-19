package jupiter

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetSwapQuote(t *testing.T) {
	quote := SwapQuoteResponse{
		InputMint:            "SOL111",
		InAmount:             "1000000000",
		OutputMint:           "USDC111",
		OutAmount:            "15025000",
		OtherAmountThreshold: "15000000",
		SwapMode:             "ExactIn",
		SlippageBps:          50,
		PriceImpactPct:       "0.01",
		RoutePlan: []RoutePlanStep{
			{
				SwapInfo: SwapInfo{
					AmmKey:     "amm123",
					Label:      "Raydium",
					InputMint:  "SOL111",
					OutputMint: "USDC111",
					InAmount:   "1000000000",
					OutAmount:  "15025000",
				},
			},
		},
	}

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/swap/v1/quote" {
			t.Errorf("expected path /swap/v1/quote, got %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("inputMint") != "SOL111" {
			t.Errorf("expected inputMint=SOL111, got %s", q.Get("inputMint"))
		}
		if q.Get("outputMint") != "USDC111" {
			t.Errorf("expected outputMint=USDC111, got %s", q.Get("outputMint"))
		}
		if q.Get("amount") != "1000000000" {
			t.Errorf("expected amount=1000000000, got %s", q.Get("amount"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(quote)
	})
	client := newTestClient(server.URL)

	result, err := client.GetSwapQuote(context.Background(), SwapQuoteParams{
		InputMint:  "SOL111",
		OutputMint: "USDC111",
		Amount:     "1000000000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.InputMint != "SOL111" {
		t.Errorf("expected InputMint SOL111, got %s", result.InputMint)
	}
	if result.OutAmount != "15025000" {
		t.Errorf("expected OutAmount 15025000, got %s", result.OutAmount)
	}
	if len(result.RoutePlan) != 1 {
		t.Fatalf("expected 1 route plan step, got %d", len(result.RoutePlan))
	}
	if result.RoutePlan[0].SwapInfo.Label != "Raydium" {
		t.Errorf("expected label Raydium, got %s", result.RoutePlan[0].SwapInfo.Label)
	}
}

func TestGetSwapQuote_AllParams(t *testing.T) {
	restrict := true

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("slippageBps") != "100" {
			t.Errorf("expected slippageBps=100, got %s", q.Get("slippageBps"))
		}
		if q.Get("swapMode") != "ExactOut" {
			t.Errorf("expected swapMode=ExactOut, got %s", q.Get("swapMode"))
		}
		if q.Get("dexes") != "Raydium" {
			t.Errorf("expected dexes=Raydium, got %s", q.Get("dexes"))
		}
		if q.Get("excludeDexes") != "Orca" {
			t.Errorf("expected excludeDexes=Orca, got %s", q.Get("excludeDexes"))
		}
		if q.Get("restrictIntermediateTokens") != "true" {
			t.Errorf("expected restrictIntermediateTokens=true, got %s", q.Get("restrictIntermediateTokens"))
		}
		if q.Get("onlyDirectRoutes") != "true" {
			t.Errorf("expected onlyDirectRoutes=true, got %s", q.Get("onlyDirectRoutes"))
		}
		if q.Get("asLegacyTransaction") != "true" {
			t.Errorf("expected asLegacyTransaction=true, got %s", q.Get("asLegacyTransaction"))
		}
		if q.Get("platformFeeBps") != "25" {
			t.Errorf("expected platformFeeBps=25, got %s", q.Get("platformFeeBps"))
		}
		if q.Get("maxAccounts") != "64" {
			t.Errorf("expected maxAccounts=64, got %s", q.Get("maxAccounts"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SwapQuoteResponse{})
	})
	client := newTestClient(server.URL)

	_, err := client.GetSwapQuote(context.Background(), SwapQuoteParams{
		InputMint:                  "SOL",
		OutputMint:                 "USDC",
		Amount:                     "1000",
		SlippageBps:                100,
		SwapMode:                   "ExactOut",
		Dexes:                      "Raydium",
		ExcludeDexes:               "Orca",
		RestrictIntermediateTokens: &restrict,
		OnlyDirectRoutes:           true,
		AsLegacyTransaction:        true,
		PlatformFeeBps:             25,
		MaxAccounts:                64,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSwapQuote_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "missing params"))
	client := newTestClient(server.URL)

	_, err := client.GetSwapQuote(context.Background(), SwapQuoteParams{})
	if err == nil {
		t.Fatal("expected error")
	}
}

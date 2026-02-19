package jupiter

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/trigger/v1/createOrder" {
			t.Errorf("expected path /trigger/v1/createOrder, got %s", r.URL.Path)
		}

		body, _ := io.ReadAll(r.Body)
		var req CreateOrderRequest
		json.Unmarshal(body, &req)
		if req.InputMint != "SOL" {
			t.Errorf("expected InputMint SOL, got %s", req.InputMint)
		}
		if req.OutputMint != "USDC" {
			t.Errorf("expected OutputMint USDC, got %s", req.OutputMint)
		}
		if req.Maker != "maker1" {
			t.Errorf("expected Maker maker1, got %s", req.Maker)
		}
		if req.Params.MakingAmount != "1000000" {
			t.Errorf("expected MakingAmount 1000000, got %s", req.Params.MakingAmount)
		}
		if req.Params.TakingAmount != "150000" {
			t.Errorf("expected TakingAmount 150000, got %s", req.Params.TakingAmount)
		}

		resp := CreateOrderResponse{
			Order:       "order123",
			Transaction: "tx456",
			RequestID:   "req789",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})
	client := newTestClient(server.URL)

	result, err := client.CreateOrder(context.Background(), CreateOrderRequest{
		InputMint:  "SOL",
		OutputMint: "USDC",
		Maker:      "maker1",
		Payer:      "payer1",
		Params: CreateOrderParams{
			MakingAmount: "1000000",
			TakingAmount: "150000",
		},
		ComputeUnitPrice: "1000",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Order != "order123" {
		t.Errorf("expected order order123, got %s", result.Order)
	}
	if result.Transaction != "tx456" {
		t.Errorf("expected transaction tx456, got %s", result.Transaction)
	}
	if result.RequestID != "req789" {
		t.Errorf("expected requestId req789, got %s", result.RequestID)
	}
}

func TestCreateOrder_WithOptionalParams(t *testing.T) {
	wrapSol := true

	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req CreateOrderRequest
		json.Unmarshal(body, &req)

		if req.Params.SlippageBps != "50" {
			t.Errorf("expected SlippageBps 50, got %v", req.Params.SlippageBps)
		}
		if req.Params.FeeBps != "10" {
			t.Errorf("expected FeeBps 10, got %v", req.Params.FeeBps)
		}
		if req.Params.ExpiredAt != "1700000000" {
			t.Errorf("expected ExpiredAt 1700000000, got %v", req.Params.ExpiredAt)
		}
		if req.WrapAndUnwrapSol == nil || !*req.WrapAndUnwrapSol {
			t.Errorf("expected WrapAndUnwrapSol true, got %v", req.WrapAndUnwrapSol)
		}
		if req.FeeAccount != "feeAcct" {
			t.Errorf("expected FeeAccount feeAcct, got %s", req.FeeAccount)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateOrderResponse{})
	})
	client := newTestClient(server.URL)

	_, err := client.CreateOrder(context.Background(), CreateOrderRequest{
		InputMint:  "SOL",
		OutputMint: "USDC",
		Maker:      "maker1",
		Payer:      "payer1",
		Params: CreateOrderParams{
			MakingAmount: "1000000",
			TakingAmount: "150000",
			SlippageBps:  "50",
			FeeBps:       "10",
			ExpiredAt:    "1700000000",
		},
		ComputeUnitPrice: "1000",
		FeeAccount:       "feeAcct",
		WrapAndUnwrapSol: &wrapSol,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateOrder_Error(t *testing.T) {
	server := newTestServer(t, errorHandler(http.StatusBadRequest, "invalid order"))
	client := newTestClient(server.URL)

	_, err := client.CreateOrder(context.Background(), CreateOrderRequest{})
	if err == nil {
		t.Fatal("expected error")
	}
}

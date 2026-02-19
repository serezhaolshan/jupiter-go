package jupiter

import (
	"context"
)

type CreateOrderParams struct {
	MakingAmount string `json:"makingAmount"`
	TakingAmount string `json:"takingAmount"`
	SlippageBps  *int   `json:"slippageBps,omitempty"`
	ExpiredAt    *int64 `json:"expiredAt,omitempty"`
	FeeBps       *int   `json:"feeBps,omitempty"`
}

type CreateOrderRequest struct {
	InputMint        string            `json:"inputMint"`
	OutputMint       string            `json:"outputMint"`
	Maker            string            `json:"maker"`
	Payer            string            `json:"payer"`
	Params           CreateOrderParams `json:"params"`
	ComputeUnitPrice string            `json:"computeUnitPrice"`
	FeeAccount       string            `json:"feeAccount,omitempty"`
	WrapAndUnwrapSol *bool             `json:"wrapAndUnwrapSol,omitempty"`
}

type CreateOrderResponse struct {
	Order       string `json:"order"`
	Transaction string `json:"transaction"`
	RequestID   string `json:"requestId"`
}

func (c *Client) CreateOrder(ctx context.Context, body CreateOrderRequest) (*CreateOrderResponse, error) {
	request, err := NewPostRequest(c.Url("/trigger/v1/createOrder"), body)
	if err != nil {
		return nil, err
	}
	var response CreateOrderResponse
	_, err = c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

package jupiter

import (
	"context"
)

type CancelOrderRequest struct {
	Maker            string `json:"maker"`
	Order            string `json:"order"`
	ComputeUnitPrice string `json:"computeUnitPrice"`
}

type CancelOrderResponse struct {
	Transaction string `json:"transaction"`
	RequestID   string `json:"requestId"`
}

func (c *Client) CancelOrder(ctx context.Context, body CancelOrderRequest) (*CancelOrderResponse, error) {
	request, err := NewPostRequest(c.Url("/trigger/v1/cancelOrder"), body)
	if err != nil {
		return nil, err
	}
	var response CancelOrderResponse
	_, err = c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

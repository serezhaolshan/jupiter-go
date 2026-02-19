package jupiter

import (
	"context"
	"fmt"
	"net/url"
)

type GetTriggerOrdersParams struct {
	User        string
	OrderStatus string
	InputMint   string
	OutputMint  string
	Page        int
}

type TriggerOrder struct {
	ID         string `json:"id"`
	User       string `json:"user"`
	InputMint  string `json:"inputMint"`
	OutputMint string `json:"outputMint"`
	Status     string `json:"status"`
}

type GetTriggerOrdersResponse struct {
	Orders      []TriggerOrder `json:"orders"`
	HasMoreData bool           `json:"hasMoreData"`
	Page        int            `json:"page"`
}

func (c *Client) GetTriggerOrders(ctx context.Context, params GetTriggerOrdersParams) (*GetTriggerOrdersResponse, error) {
	queryParams := url.Values{}

	queryParams.Set("user", params.User)
	queryParams.Set("orderStatus", params.OrderStatus)

	if params.InputMint != "" {
		queryParams.Set("inputMint", params.InputMint)
	}
	if params.OutputMint != "" {
		queryParams.Set("outputMint", params.OutputMint)
	}
	if params.Page > 0 {
		queryParams.Set("page", fmt.Sprintf("%d", params.Page))
	}

	request := NewRequest(c.Url("/trigger/v1/getTriggerOrders"), queryParams)
	var response GetTriggerOrdersResponse
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

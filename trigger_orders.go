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
	UserPubkey              string   `json:"userPubkey"`
	OrderKey                string   `json:"orderKey"`
	InputMint               string   `json:"inputMint"`
	OutputMint              string   `json:"outputMint"`
	MakingAmount            string   `json:"makingAmount"`
	TakingAmount            string   `json:"takingAmount"`
	RemainingMakingAmount   string   `json:"remainingMakingAmount"`
	RemainingTakingAmount   string   `json:"remainingTakingAmount"`
	RawMakingAmount         string   `json:"rawMakingAmount"`
	RawTakingAmount         string   `json:"rawTakingAmount"`
	RawRemainingMakingAmount string  `json:"rawRemainingMakingAmount"`
	RawRemainingTakingAmount string  `json:"rawRemainingTakingAmount"`
	SlippageBps             string   `json:"slippageBps"`
	ExpiredAt               *string  `json:"expiredAt"`
	CreatedAt               string   `json:"createdAt"`
	UpdatedAt               string   `json:"updatedAt"`
	Status                  string   `json:"status"`
	OpenTx                  string   `json:"openTx"`
	CloseTx                 *string  `json:"closeTx"`
	ProgramVersion          string   `json:"programVersion"`
	Trades                  []any    `json:"trades"`
}

type GetTriggerOrdersResponse struct {
	User        string         `json:"user"`
	OrderStatus string         `json:"orderStatus"`
	Orders      []TriggerOrder `json:"orders"`
	TotalPages  int            `json:"totalPages"`
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

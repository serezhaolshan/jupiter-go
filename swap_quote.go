package jupiter

import (
	"context"
	"fmt"
	"net/url"
)

type SwapQuoteParams struct {
	InputMint                  string
	OutputMint                 string
	Amount                     string
	SlippageBps                int
	SwapMode                   string
	Dexes                      string
	ExcludeDexes               string
	RestrictIntermediateTokens *bool
	OnlyDirectRoutes           bool
	AsLegacyTransaction        bool
	PlatformFeeBps             int
	MaxAccounts                int
}

type SwapQuoteResponse struct {
	InputMint            string          `json:"inputMint"`
	InAmount             string          `json:"inAmount"`
	OutputMint           string          `json:"outputMint"`
	OutAmount            string          `json:"outAmount"`
	OtherAmountThreshold string          `json:"otherAmountThreshold"`
	SwapMode             string          `json:"swapMode"`
	SlippageBps          int             `json:"slippageBps"`
	PriceImpactPct       string          `json:"priceImpactPct"`
	RoutePlan            []RoutePlanStep `json:"routePlan"`
	PlatformFee          *PlatformFee    `json:"platformFee,omitempty"`
	ContextSlot          *int64          `json:"contextSlot,omitempty"`
	TimeTaken            *float64        `json:"timeTaken,omitempty"`
}

func (c *Client) GetSwapQuote(ctx context.Context, params SwapQuoteParams) (*SwapQuoteResponse, error) {
	queryParams := url.Values{}

	queryParams.Set("inputMint", params.InputMint)
	queryParams.Set("outputMint", params.OutputMint)
	queryParams.Set("amount", params.Amount)

	if params.SlippageBps > 0 {
		queryParams.Set("slippageBps", fmt.Sprintf("%d", params.SlippageBps))
	}
	if params.SwapMode != "" {
		queryParams.Set("swapMode", params.SwapMode)
	}
	if params.Dexes != "" {
		queryParams.Set("dexes", params.Dexes)
	}
	if params.ExcludeDexes != "" {
		queryParams.Set("excludeDexes", params.ExcludeDexes)
	}
	if params.RestrictIntermediateTokens != nil {
		queryParams.Set("restrictIntermediateTokens", fmt.Sprintf("%t", *params.RestrictIntermediateTokens))
	}
	if params.OnlyDirectRoutes {
		queryParams.Set("onlyDirectRoutes", "true")
	}
	if params.AsLegacyTransaction {
		queryParams.Set("asLegacyTransaction", "true")
	}
	if params.PlatformFeeBps > 0 {
		queryParams.Set("platformFeeBps", fmt.Sprintf("%d", params.PlatformFeeBps))
	}
	if params.MaxAccounts > 0 {
		queryParams.Set("maxAccounts", fmt.Sprintf("%d", params.MaxAccounts))
	}

	request := NewRequest(c.Url("/swap/v1/quote"), queryParams)
	var response SwapQuoteResponse
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

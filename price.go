package jupiter

import (
	"context"
	"net/url"
)

type PriceV3Entry struct {
	USDPrice       float64 `json:"usdPrice"`
	BlockID        *int64  `json:"blockId,omitempty"`
	Decimals       *int    `json:"decimals,omitempty"`
	PriceChange24h *float64 `json:"priceChange24h,omitempty"`
}

type PriceV3Response map[string]PriceV3Entry

func (c *Client) GetPrices(ctx context.Context, ids string) (PriceV3Response, error) {
	queryParams := url.Values{}
	queryParams.Set("ids", ids)

	request := NewRequest(c.Url("/price/v3"), queryParams)
	var response PriceV3Response
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

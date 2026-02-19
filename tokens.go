package jupiter

import (
	"context"
	"fmt"
	"net/url"
)

type TokenV2 struct {
	ID                string      `json:"id"`
	Name              string      `json:"name"`
	Symbol            string      `json:"symbol"`
	Icon              string      `json:"icon,omitempty"`
	Decimals          int         `json:"decimals"`
	CircSupply        *float64    `json:"circSupply,omitempty"`
	TotalSupply       *float64    `json:"totalSupply,omitempty"`
	TokenProgram      string      `json:"tokenProgram,omitempty"`
	FirstPool         *FirstPool  `json:"firstPool,omitempty"`
	HolderCount       *int        `json:"holderCount,omitempty"`
	Audit             *Audit      `json:"audit,omitempty"`
	OrganicScore      *float64    `json:"organicScore,omitempty"`
	OrganicScoreLabel string      `json:"organicScoreLabel,omitempty"`
	IsVerified        *bool       `json:"isVerified,omitempty"`
	Cexes             []string    `json:"cexes,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	FDV               *float64    `json:"fdv,omitempty"`
	MCap              *float64    `json:"mcap,omitempty"`
	USDPrice          *float64    `json:"usdPrice,omitempty"`
	PriceBlockID      *int64      `json:"priceBlockId,omitempty"`
	Liquidity         *float64    `json:"liquidity,omitempty"`
	Stats5m           *TokenStats `json:"stats5m,omitempty"`
	Stats1h           *TokenStats `json:"stats1h,omitempty"`
	Stats6h           *TokenStats `json:"stats6h,omitempty"`
	Stats24h          *TokenStats `json:"stats24h,omitempty"`
	UpdatedAt         string      `json:"updatedAt,omitempty"`
}

type GetTokensParams struct {
	SortBy   string
	Interval string
	Limit    int
}

func (c *Client) GetTokens(ctx context.Context, params GetTokensParams) ([]TokenV2, error) {
	queryParams := url.Values{}

	if params.Interval != "" {
		queryParams.Set("interval", params.Interval)
	}
	if params.Limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", params.Limit))
	}

	endpoint := fmt.Sprintf("/tokens/v2/%s", params.SortBy)
	request := NewRequest(c.Url(endpoint), queryParams)
	var response []TokenV2
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

type SearchTokensParams struct {
	Query string
}

func (c *Client) SearchTokens(ctx context.Context, params SearchTokensParams) ([]TokenV2, error) {
	queryParams := url.Values{}
	queryParams.Set("query", params.Query)

	request := NewRequest(c.Url("/tokens/v2/search"), queryParams)
	var response []TokenV2
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

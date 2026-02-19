package jupiter

import (
	"context"
	"net/url"
)

type Router struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon,omitempty"`
}

type RoutersResponse []Router

func (c *Client) GetRouters(ctx context.Context) (RoutersResponse, error) {
	request := NewRequest(c.Url("/ultra/v1/order/routers"), url.Values{})
	var response RoutersResponse
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

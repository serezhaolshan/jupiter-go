package jupiter

import (
	"context"
	"net/url"
)

type RoutersResponse []string

func (c *Client) GetRouters(ctx context.Context) (RoutersResponse, error) {
	request := NewRequest(c.Url("/ultra/v1/order/routers"), url.Values{})
	var response RoutersResponse
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

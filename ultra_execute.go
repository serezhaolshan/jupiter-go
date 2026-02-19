package jupiter

import (
	"context"
)

func (c *Client) ExecuteUltra(ctx context.Context, body ExecuteRequest) (*ExecuteResponse, error) {
	request, err := NewPostRequest(c.Url("/ultra/v1/execute"), body)
	if err != nil {
		return nil, err
	}
	var response ExecuteResponse
	_, err = c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

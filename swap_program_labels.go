package jupiter

import (
	"context"
	"net/url"
)

type ProgramIDToLabelResponse map[string]string

func (c *Client) GetProgramIDToLabel(ctx context.Context) (ProgramIDToLabelResponse, error) {
	request := NewRequest(c.Url("/swap/v1/program-id-to-label"), url.Values{})
	var response ProgramIDToLabelResponse
	_, err := c.doCall(ctx, request, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

package sourcecraft

import (
	"context"
)

func (c *Client) GetCurrentUser(ctx context.Context) (*UserProfile, *Response, error) {
	userResp := UserProfile{}
	resp, err := c.getParsedResponse(ctx, "GET", "/user", nil, nil, &userResp)
	return &userResp, resp, err
}

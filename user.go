package lunchmoney

import (
	"context"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
)

func (c *Client) FetchMe(ctx context.Context) (*models.User, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "me", nil, nil)
	if err != nil {
		return nil, err
	}

	var me models.User
	err = c.Do(ctx, req, &me)

	return &me, err
}

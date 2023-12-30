package lunchmoney

import (
	"context"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

func (c *Client) FetchMe(ctx context.Context) (*models.User, error) {
	req, err := http.NewRequest(http.MethodGet, constants.Path.Me, nil)
	if err != nil {
		return nil, err
	}

	var me models.User
	err = c.Do(ctx, req, &me, nil)

	return &me, err
}

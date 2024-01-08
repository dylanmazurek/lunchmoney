package lunchmoney

import (
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

func (c *Client) FetchMe() (*models.User, error) {
	req, err := c.NewRequest(http.MethodGet, constants.Path.Me, nil, nil)
	if err != nil {
		return nil, err
	}

	var me models.User
	err = c.Do(req, &me)

	return &me, err
}

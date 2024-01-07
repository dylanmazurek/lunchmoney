package lunchmoney

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

func (c *Client) ListCategory() (*[]models.Category, error) {
	urlString := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Categories)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(http.MethodGet, requestUrl.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	var categories models.CategoryResponse
	err = c.Do(req, &categories)

	return &categories.Categories, err
}

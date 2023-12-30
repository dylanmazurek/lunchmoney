package lunchmoney

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

func (c *Client) ListCategory(ctx context.Context) (*[]models.Category, error) {
	urlString := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Categories)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	var categories models.CategoryResponse
	err = c.Do(ctx, req, &categories, nil)

	return &categories.Categories, err
}

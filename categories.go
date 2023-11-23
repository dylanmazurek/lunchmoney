package lunchmoney

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
)

// GetCategories gets all categories
func (c *Client) GetCategories(ctx context.Context) (categories *[]models.Category, err error) {
	path := "/v1/categories"

	reqOptions := models.Request{
		Method: http.MethodGet,
		Path:   path,
	}

	response, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get categories: %w", err)
	}

	return response.Categories, nil
}

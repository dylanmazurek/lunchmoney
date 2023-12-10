package lunchmoney

import (
	"context"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
)

func (c *Client) ListCategory(ctx context.Context) (*[]models.Category, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "categories", nil, nil)
	if err != nil {
		return nil, err
	}

	var categories []models.Category
	err = c.Do(ctx, req, &categories)

	return &categories, nil
}

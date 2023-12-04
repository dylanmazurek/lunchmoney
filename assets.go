package lunchmoney

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/dylanmazurek/lunchmoney/models"
)

// FetchAsset gets a single asset based on id
func (c *Client) FetchAsset(ctx context.Context, assetId int64) (asset *models.Asset, err error) {
	assets, err := c.ListAssets(ctx)
	if err != nil {
		return nil, err
	}

	assetIdx := slices.IndexFunc(*assets, func(asset models.Asset) bool { return asset.ID == assetId })
	if assetIdx == -1 {
		return nil, nil
	}

	currentAsset := (*assets)[assetIdx]

	return &currentAsset, err
}

// ListAssets gets all assets filtered by the filters.
func (c *Client) ListAssets(ctx context.Context) (assets *[]models.Asset, err error) {
	path := "/v1/assets"

	reqOptions := models.Request{
		Method: http.MethodGet,
		Path:   path,
	}

	body, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("list assets: %w", err)
	}

	if body == nil {
		return nil, nil
	}

	return body.Assets, err
}

// UpdateAsset updates a transaction by id.
func (c *Client) UpdateAsset(ctx context.Context, assetId int64, asset models.Asset) (*models.Asset, error) {
	path := fmt.Sprintf("/v1/assets/%d", assetId)

	reqOptions := models.Request{
		Method:  http.MethodPut,
		Path:    path,
		ReqBody: asset,
	}

	resp, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("update asset: %w", err)
	}

	wasUpdated := resp.Updated
	if wasUpdated == nil || !*wasUpdated {
		return nil, fmt.Errorf("asset was not updated")
	}

	return resp.Item, err
}

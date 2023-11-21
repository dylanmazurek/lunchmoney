package lunchmoney

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
)

// // GetAssets gets all assets filtered by the filters.
func (c *Client) GetAssets(ctx context.Context) (assets *[]models.Asset, err error) {
	path := "/v1/assets/"

	reqOptions := models.RequestOptions{
		Method: http.MethodGet,
		Path:   path,
	}

	body, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get assets: %w", err)
	}

	return body.Assets, err
}

// UpdateAsset updates a transaction by id.
func (c *Client) UpdateAsset(ctx context.Context, assetId int64, reqBody *models.AssetUpdateRequest) (asset *models.Asset, err error) {
	path := fmt.Sprintf("/v1/assets/%d", assetId)

	reqOptions := models.RequestOptions{
		Method:  http.MethodPut,
		Path:    path,
		ReqBody: reqBody,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("update asset: %w", err)
	}

	if !respBody.Updated {
		return nil, fmt.Errorf("asset was not updated")
	}

	return respBody.Asset, err
}

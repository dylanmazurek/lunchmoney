package lunchmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

func (c *Client) FetchAsset(ctx context.Context, assetId int64) (*models.Asset, error) {
	assets, err := c.ListAsset(ctx)
	if err != nil {
		return nil, err
	}

	assetIdx := slices.IndexFunc(*assets, func(asset models.Asset) bool { return *asset.AssetID == assetId })
	if assetIdx == -1 {
		return nil, nil
	}

	asset := (*assets)[assetIdx]

	return &asset, err
}

func (c *Client) ListAsset(ctx context.Context) (*[]models.Asset, error) {
	urlString := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Assets)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	var assets models.AssetResponse
	err = c.Do(ctx, req, &assets, nil)

	return &assets.Assets, err
}

func (c *Client) UpdateAsset(ctx context.Context, id int64, asset models.Asset) (*models.Asset, error) {
	assetJson, err := json.Marshal(&asset)
	if err != nil {
		return nil, err
	}

	urlString := fmt.Sprintf("%s/%s/%d", constants.Config.APIBaseURL, constants.Path.Assets, id)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, requestUrl.String(), bytes.NewReader(assetJson))
	if err != nil {
		return nil, err
	}

	var updatedAsset models.Asset
	err = c.Do(ctx, req, &updatedAsset, nil)

	if err != nil {
		return nil, err
	}

	if updatedAsset.Error != nil {
		return nil, errors.New(*updatedAsset.Error)
	}

	return &updatedAsset, err
}

package lunchmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

func (c *Client) FetchAsset(assetId int64) (*models.Asset, error) {
	assets, err := c.ListAsset()
	if err != nil {
		return nil, err
	}

	assetIdx := slices.IndexFunc(*assets, func(asset models.Asset) bool { return asset.AssetID == &assetId })
	if assetIdx == -1 {
		return nil, nil
	}

	asset := (*assets)[assetIdx]

	return &asset, err
}

func (c *Client) ListAsset() (*[]models.Asset, error) {
	urlString := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Assets)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest(http.MethodGet, requestUrl.String(), nil, nil)
	if err != nil {
		return nil, err
	}

	var assets models.AssetResponse
	err = c.Do(req, &assets)

	return &assets.Assets, err
}

func (c *Client) UpdateAsset(id int64, asset *models.Asset) (*models.Asset, error) {
	assetJson, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	requestPath := fmt.Sprintf("%s/%d", constants.Path.Assets, id)

	req, err := c.NewRequest(http.MethodPut, requestPath, bytes.NewReader(assetJson), nil)
	if err != nil {
		return nil, err
	}

	var updatedAsset models.Asset
	err = c.Do(req, updatedAsset)

	if err != nil {
		return nil, err
	}

	if updatedAsset.Error != nil {
		return nil, errors.New(*updatedAsset.Error)
	}

	return &updatedAsset, err
}

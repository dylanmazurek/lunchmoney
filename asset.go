package lunchmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/dylanmazurek/lunchmoney/models"
)

func (c *Client) FetchAsset(ctx context.Context, assetId int64) (*models.Asset, error) {
	assets, err := c.ListAsset(ctx)
	if err != nil {
		return nil, err
	}

	assetIdx := slices.IndexFunc(assets.Assets, func(asset models.Asset) bool { return *asset.AssetID == assetId })
	if assetIdx == -1 {
		return nil, nil
	}

	asset := (assets.Assets)[assetIdx]

	return &asset, err
}

func (c *Client) ListAsset(ctx context.Context) (*models.AssetList, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, "assets", nil, nil)
	if err != nil {
		return nil, err
	}

	var assets models.AssetList
	err = c.Do(ctx, req, &assets)

	return &assets, err
}

func (c *Client) UpdateAsset(ctx context.Context, id int64, asset models.Asset) (*models.Asset, error) {
	b, err := json.Marshal(&asset)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("assets/%d", id)
	req, err := c.NewRequest(ctx, http.MethodPut, path, bytes.NewReader(b), nil)
	if err != nil {
		return nil, err
	}

	var updatedAsset models.Asset
	err = c.Do(ctx, req, &updatedAsset)

	if updatedAsset.Error != nil {
		err = fmt.Errorf("%s", *updatedAsset.Error)
	}

	return &updatedAsset, err
}

func (c *Client) UpdateAssetFromJSON(ctx context.Context, assetJson []byte) (*models.Asset, error) {
	asset := &models.Asset{}
	err := json.Unmarshal(assetJson, &asset)

	if asset.AssetID == nil {
		return nil, errors.New("no asset id set")
	}

	updatedAsset, err := c.UpdateAsset(ctx, *asset.AssetID, *asset)

	return updatedAsset, err
}

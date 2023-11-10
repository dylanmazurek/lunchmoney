package lunchmoney

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dylanmazurek/lunchmoney/models"
)

// // AssetsResponse is a response to an asset lookup.
// type AssetsResponse struct {
// 	Assets []*Asset `json:"assets"`
// }

// // Asset is a single LM asset.
// type Asset struct {
// 	ID              int64     `json:"id"`
// 	TypeName        string    `json:"type_name"`
// 	SubtypeName     string    `json:"subtype_name"`
// 	Name            string    `json:"name"`
// 	Balance         string    `json:"balance"`
// 	BalanceAsOf     time.Time `json:"balance_as_of"`
// 	Currency        string    `json:"currency"`
// 	Status          string    `json:"status"`
// 	InstitutionName string    `json:"institution_name"`
// 	CreatedAt       time.Time `json:"created_at"`
// }

// // GetAssets gets all assets filtered by the filters.
// func (c *Client) GetAssets(ctx context.Context) ([]*Asset, error) {
// 	validate := validator.New()
// 	options := map[string]string{}

// 	body, err := c.Get(ctx, "/v1/assets", options)
// 	if err != nil {
// 		return nil, fmt.Errorf("get assets: %w", err)
// 	}

// 	resp := &AssetsResponse{}
// 	if err := json.NewDecoder(body).Decode(resp); err != nil {
// 		return nil, fmt.Errorf("decode response: %w", err)
// 	}

// 	if err := validate.Struct(resp); err != nil {
// 		return nil, err
// 	}

// 	return resp.Assets, nil
// }

// UpdateAsset updates a transaction by id.
func (c *Client) UpdateAsset(ctx context.Context, assetId int64, body *models.AssetUpdateRequest) (resp *models.AssetsResponse, err error) {
	path := fmt.Sprintf("/v1/assets/%d", assetId)

	reqOptions := models.RequestOptions{
		Method:      http.MethodPut,
		Path:        path,
		QueryValues: nil,
		ReqBody:     body,
	}

	resp, err = Request[models.AssetsResponse](ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("update asset: %w", err)
	}

	if resp.Errors != nil {
		return nil, errors.New(strings.Join(*resp.Errors, ", "))
	}

	return resp, nil
}

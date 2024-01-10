package shared

import (
	"encoding/json"
	"time"
)

type Asset struct {
	ExternalAssetID string      `json:"externalAssetId,omitempty"`
	ExternalUserID  string      `json:"externalUserId,omitempty"`
	Balance         json.Number `json:"balance"`
	Currency        string      `json:"currency"`
	BalanceAsOf     *time.Time  `json:"balance_as_of"`

	AssetID *int64 `json:"assetId,omitempty"`
}

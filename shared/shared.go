package shared

import (
	"encoding/json"
	"time"
)

type Asset struct {
	ExternalAssetID *string     `json:"externalAssetId,omitempty"`
	Balance         json.Number `json:"balance"`
	Currency        *string     `json:"currency"`
	BalanceAsOf     *time.Time  `json:"balance_as_of"`

	AssetID int64 `json:"assetID"`
}

type Transaction struct {
	ExternalTransactionID string    `json:"externalTransactionId"`
	ExternalAssetID       string    `json:"externalAssetId"`
	Datetime              time.Time `json:"datetime"`
	Description           string    `json:"description"`
	Status                string    `json:"status"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`

	AssetID int64 `json:"assetID"`
}

package shared

import (
	"encoding/json"
	"time"
)

type Transaction struct {
	ExternalTransactionID string      `json:"externalTransactionId" bson:"externalTransactionId"`
	ExternalAssetID       string      `json:"externalAssetId" bson:"externalAssetId"`
	Datetime              time.Time   `json:"datetime"`
	Description           string      `json:"description"`
	Status                string      `json:"status"`
	Amount                json.Number `json:"amount"`
	Currency              string      `json:"currency"`

	AssetID *int64 `json:"assetId"`
}

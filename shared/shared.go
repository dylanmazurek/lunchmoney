package shared

import (
	"encoding/json"
	"time"
)

// pub sub

type Job struct {
	Action   string     `json:"action"`
	AssetID  *int       `json:"assetId,omitempty"`
	DateFrom *time.Time `json:"dateFrom,omitempty"`
	DateTo   *time.Time `json:"dateTo,omitempty"`
}

// database

type Account struct {
	AssetSource     string `bson:"assetSource"`
	ExternalAssetID string `bson:"externalAssetId"`
	ExternalUserID  string `bson:"externalUserId"`
	AssetName       string `bson:"assetName"`
	Currency        string `bson:"currency"`
	AssetID         int64  `bson:"assetId"`
}

type Asset struct {
	ExternalAssetID string      `json:"externalAssetId,omitempty"`
	ExternalUserID  string      `json:"externalUserId,omitempty"`
	Balance         json.Number `json:"balance"`
	Currency        string      `json:"currency"`
	BalanceAsOf     *time.Time  `json:"balance_as_of"`

	AssetID int64 `json:"assetId"`
}

type Transaction struct {
	ExternalTransactionID string    `json:"externalTransactionId"`
	ExternalAssetID       string    `json:"externalAssetId"`
	Datetime              time.Time `json:"datetime"`
	Description           string    `json:"description"`
	Status                string    `json:"status"`
	Amount                float64   `json:"amount"`
	Currency              string    `json:"currency"`

	AssetID int64 `json:"assetId"`
}

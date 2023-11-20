package models

import "time"

type Asset struct {
	ID                  int       `json:"id"`
	TypeName            string    `json:"type_name"`
	SubtypeName         string    `json:"subtype_name"`
	Name                string    `json:"name"`
	Balance             string    `json:"balance"`
	BalanceAsOf         time.Time `json:"balance_as_of"`
	Currency            string    `json:"currency"`
	InstitutionName     string    `json:"institution_name"`
	ExcludeTransactions bool      `json:"exclude_transactions"`
	CreatedAt           string    `json:"created_at"`
}

// AssetUpdateRequest
type AssetUpdateRequest struct {
	Balance string `json:"balance"`
}

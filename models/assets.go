package models

import (
	"encoding/json"
	"time"

	"github.com/Rhymond/go-money"
)

type Asset struct {
	ID                  int64  `json:"id"`
	TypeName            string `json:"type_name"`
	SubtypeName         string `json:"subtype_name"`
	Name                string `json:"name"`
	DisplayName         string `json:"display_name"`
	InstitutionName     string `json:"institution_name"`
	ExcludeTransactions bool   `json:"exclude_transactions"`
	CreatedAt           string `json:"created_at"`

	Balance     money.Money `json:"-"`
	BalanceAsOf time.Time   `json:"balance_as_of"`
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	type Alias Asset
	marshaledJSON, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})

	return marshaledJSON, err
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	type Alias Asset
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

package models

import (
	"encoding/json"
	"math"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

type AssetList struct {
	Assets []Asset `json:"assets"`
}

func (a *AssetList) UnmarshalJSON(data []byte) error {
	type Alias AssetList
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

type Asset struct {
	ID                  *int64  `json:"id,omitempty"`
	TypeName            *string `json:"type_name,omitempty"`
	SubtypeName         *string `json:"subtype_name,omitempty"`
	Name                *string `json:"name,omitempty"`
	DisplayName         *string `json:"display_name,omitempty"`
	InstitutionName     string  `json:"institution_name,omitempty"`
	ExcludeTransactions *bool   `json:"exclude_transactions,omitempty"`
	CreatedAt           *string `json:"created_at,omitempty"`

	Balance     money.Money    `json:"-"`
	Currency    money.Currency `json:"-"`
	BalanceAsOf *time.Time     `json:"balance_as_of,omitempty"`
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	type Alias Asset
	marshaledJSON, err := json.Marshal(&struct {
		*Alias
		BalanceRaw  json.Number `json:"balance"`
		CurrencyRaw string      `json:"currency"`
	}{
		Alias:       (*Alias)(a),
		BalanceRaw:  json.Number(a.Balance.Amount()),
		CurrencyRaw: strings.ToLower(a.Currency.Code),
	})

	return marshaledJSON, err
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	type Alias Asset
	aux := &struct {
		*Alias
		BalanceRaw  json.Number `json:"balance"`
		CurrencyRaw string      `json:"currency"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	currency := money.GetCurrency(strings.ToUpper(aux.CurrencyRaw))
	a.Currency = *currency

	balanceFloat, _ := aux.BalanceRaw.Float64()
	balanceInCents := int64(math.Round(balanceFloat))
	balance := money.New(balanceInCents, a.Currency.Code)
	a.Balance = *balance

	return nil
}

package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

type AssetResponse struct {
	Assets []Asset `json:"assets"`
}

func (a *AssetResponse) UnmarshalJSON(data []byte) error {
	type Alias AssetResponse
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
	Error *string `json:"error,omitempty"`

	AssetID             *int64  `json:"-"`
	TypeName            *string `json:"type_name,omitempty"`
	SubtypeName         *string `json:"subtype_name,omitempty"`
	Name                *string `json:"name,omitempty"`
	DisplayName         *string `json:"display_name,omitempty"`
	InstitutionName     string  `json:"institution_name,omitempty"`
	ExcludeTransactions *bool   `json:"exclude_transactions,omitempty"`
	CreatedAt           *string `json:"created_at,omitempty"`

	Balance     money.Money `json:"-"`
	BalanceAsOf *time.Time  `json:"-"`
}

func (a *Asset) MarshalJSON() ([]byte, error) {
	type Alias Asset
	marshaledJSON, err := json.Marshal(&struct {
		*Alias
		BalanceRaw     string `json:"balance"`
		CurrencyRaw    string `json:"currency"`
		BalanceAsOfRaw string `json:"balance_as_of"`
	}{
		Alias:          (*Alias)(a),
		BalanceRaw:     fmt.Sprintf("%.2f", a.Balance.AsMajorUnits()),
		CurrencyRaw:    strings.ToLower(a.Balance.Currency().Code),
		BalanceAsOfRaw: a.BalanceAsOf.UTC().Format(time.RFC3339),
	})

	return marshaledJSON, err
}

func (a *Asset) UnmarshalJSON(data []byte) error {
	type Alias Asset
	aux := &struct {
		*Alias
		BalanceRaw     string `json:"balance"`
		CurrencyRaw    string `json:"currency"`
		BalanceAsOfRaw string `json:"balance_as_of"`
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.BalanceRaw != "" {
		balanceFloat, err := strconv.ParseFloat(aux.BalanceRaw, 64)
		if err != nil {
			return err
		}

		currency := strings.ToUpper(aux.CurrencyRaw)
		balance := money.NewFromFloat(balanceFloat, currency)
		a.Balance = *balance
	}

	if aux.BalanceAsOfRaw != "" {
		balanceAsOf, err := time.Parse(time.RFC3339, aux.BalanceAsOfRaw)
		if err != nil {
			return err
		}

		a.BalanceAsOf = &balanceAsOf
	}

	return nil
}

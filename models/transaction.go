package models

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Rhymond/go-money"
	"github.com/dylanmazurek/lunchmoney/util/time"
)

type TransactionFilter struct {
	AssetID int `url:"asset_id"`

	TagID           *int64     `url:"tag_id,omitempty"`
	RecurringID     *int64     `url:"recurring_id,omitempty"`
	CategoryID      *int64     `url:"category_id,omitempty"`
	StartDate       *time.Date `url:"start_date,omitempty"`
	EndDate         *time.Date `url:"end_date,omitempty"`
	DebitAsNegative bool       `url:"debit_as_negative"`

	Limit  int `url:"limit"`
	Offset int `url:"offset"`
}

type TransactionResponse struct {
	Transactions []Transaction `json:"transactions"`
}

func (a *TransactionResponse) UnmarshalJSON(data []byte) error {
	type Alias TransactionResponse
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

type Transaction struct {
	ID      *json.Number `json:"id,omitempty"`
	AssetID json.Number  `json:"asset_id"`

	Payee        string       `json:"payee"`
	Notes        *string      `json:"notes,omitempty"`
	CategoryID   *json.Number `json:"category_id,omitempty"`
	RecurringID  *json.Number `json:"recurring_id,omitempty"`
	Status       string       `json:"status"`
	IsGroup      bool         `json:"is_group,omitempty"`
	GroupID      *json.Number `json:"group_id,omitempty"`
	ParentID     *json.Number `json:"parent_id,omitempty"`
	OriginalName *string      `json:"original_name,omitempty"`
	ExternalID   *string      `json:"external_id,omitempty"`
	Type         *string      `json:"type,omitempty"`

	Amount       money.Money `json:"-"`
	Date         time.Date   `json:"-"`
	OriginalDate *time.Date  `json:"-"`
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	marshaledJSON, err := json.Marshal(&struct {
		*Alias
		AmountRaw   string `json:"amount"`
		CurrencyRaw string `json:"currency"`

		DateRaw         string  `json:"date"`
		OriginalDateRaw *string `json:"original_date,omitempty"`
	}{
		Alias:       (*Alias)(t),
		AmountRaw:   fmt.Sprintf("%.2f", t.Amount.AsMajorUnits()),
		CurrencyRaw: strings.ToLower(t.Amount.Currency().Code),

		DateRaw: t.Date.String(),
	})

	return marshaledJSON, err
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		*Alias
		AmountRaw   string `json:"amount"`
		CurrencyRaw string `json:"currency"`

		DateRaw         string  `json:"date"`
		OriginalDateRaw *string `json:"original_date,omitempty"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.AmountRaw != "" {
		amountFloat, err := strconv.ParseFloat(aux.AmountRaw, 64)
		if err != nil {
			return err
		}

		currency := strings.ToUpper(aux.CurrencyRaw)
		amount := money.NewFromFloat(amountFloat, currency)
		t.Amount = *amount
	}

	if aux.DateRaw != "" {
		t.Date = time.Parse(aux.DateRaw)
	}

	return nil
}

type InsertRequest struct {
	Transactions      []Transaction `json:"transactions"`
	ApplyRules        bool          `json:"apply_rules,omitempty"`
	SkipDuplicates    bool          `json:"skip_duplicates,omitempty"`
	CheckForRecurring bool          `json:"check_for_recurring,omitempty"`
	DebitAsNegative   bool          `json:"debit_as_negative,omitempty"`
	SkipBalanceUpdate bool          `json:"skip_balance_update,omitempty"`
}

type InsertResponse struct {
	Errors []string `json:"error,omitempty"`

	Ids []int64 `json:"ids,omitempty"`
}

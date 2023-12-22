package models

import (
	"encoding/json"
	"math"
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

type TransactionList struct {
	Transactions []Transaction `json:"transactions"`
}

func (a *TransactionList) UnmarshalJSON(data []byte) error {
	type Alias TransactionList
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
	ID           *json.Number   `json:"id,omitempty"`
	Date         time.Date      `json:"date"`
	OriginalDate *time.Date     `json:"original_date,omitempty"`
	Payee        *string        `json:"payee"`
	Amount       money.Money    `json:"-"`
	Currency     money.Currency `json:"-"`
	Notes        *string        `json:"notes,omitempty"`
	CategoryID   *json.Number   `json:"category_id,omitempty"`
	RecurringID  *json.Number   `json:"recurring_id,omitempty"`
	AssetID      json.Number    `json:"asset_id"`
	Status       *string        `json:"status,omitempty"`
	IsGroup      bool           `json:"is_group,omitempty"`
	GroupID      *json.Number   `json:"group_id,omitempty"`
	ParentID     *json.Number   `json:"parent_id,omitempty"`
	OriginalName *string        `json:"original_name,omitempty"`
	ExternalID   *string        `json:"external_id,omitempty"`
	Type         *string        `json:"type,omitempty"`
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	marshaledJSON, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	})

	return marshaledJSON, err
}

func (t *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		*Alias

		AmountRaw   json.Number `json:"amount"`
		ToBaseRaw   json.Number `json:"to_base"`
		CurrencyRaw string      `json:"currency"`
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	currency := money.GetCurrency(strings.ToUpper(aux.CurrencyRaw))
	t.Currency = *currency

	amountFloat, _ := aux.AmountRaw.Float64()
	amountInCents := int64(math.Round(amountFloat))
	amount := money.New(amountInCents, t.Currency.Code)
	t.Amount = *amount

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
	Error *string `json:"error,omitempty"`

	Ids []int64 `json:"ids,omitempty"`
}

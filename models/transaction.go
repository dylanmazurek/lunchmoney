package models

import (
	"encoding/json"

	"github.com/dylanmazurek/lunchmoney/util/money"
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

type Transaction struct {
	ID          *int64      `json:"id,omitempty"`
	Date        *time.Date  `json:"date"`
	Payee       string      `json:"payee"`
	Amount      money.Money `json:"amount"`
	Notes       *string     `json:"notes,omitempty"`
	CategoryID  *int64      `json:"category_id,omitempty"`
	RecurringID *int64      `json:"recurring_id,omitempty"`
	AssetID     int64       `json:"asset_id"`
	Status      *string     `json:"status,omitempty"`
	IsGroup     bool        `json:"is_group,omitempty"`
	GroupID     *int64      `json:"group_id,omitempty"`
	ParentID    *int64      `json:"parent_id,omitempty"`
	ExternalID  *string     `json:"external_id,omitempty"`

	Tags *[]Tag `json:"tags"`
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
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}

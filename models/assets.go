package models

import (
	"time"
)

// TransactionFilter are options to pass into the request for transactions.
type TransactionFilter struct {
	TagID           *int64     `json:"tag_id,omitempty" url:"tag_id,omitempty"`
	RecurringID     *int64     `json:"recurring_id,omitempty" url:"recurring_id,omitempty"`
	PlaidAccountID  *int64     `json:"plaid_account_id,omitempty" url:"plaid_account_id,omitempty"`
	CategoryID      *int64     `json:"category_id,omitempty" url:"category_id,omitempty"`
	AssetID         *int64     `json:"asset_id,omitempty" url:"asset_id,omitempty"`
	Offset          *int64     `json:"offset,omitempty" url:"offset,omitempty"`
	Limit           *int64     `json:"limit,omitempty" url:"limit,omitempty"`
	StartDate       *time.Time `json:"start_date,omitempty" url:"start_date,omitempty" layout:"2006-01-02"`
	EndDate         *time.Time `json:"end_date,omitempty" url:"end_date,omitempty" layout:"2006-01-02"`
	DebitAsNegative *bool      `json:"debit_as_negative,omitempty" url:"debit_as_negative,omitempty"`
}

// Transaction is a single LM transaction.
type Transaction struct {
	ID             int64            `json:"id"`
	Date           DateTimeOptional `json:"date"`
	Payee          string           `json:"payee"`
	Amount         string           `json:"amount"`
	Currency       string           `json:"currency"`
	Notes          string           `json:"notes"`
	CategoryID     int64            `json:"category_id"`
	RecurringID    int64            `json:"recurring_id"`
	AssetID        int64            `json:"asset_id"`
	PlaidAccountID int64            `json:"plaid_account_id"`
	Status         string           `json:"status"`
	IsGroup        bool             `json:"is_group"`
	GroupID        int64            `json:"group_id"`
	ParentID       int64            `json:"parent_id"`
	ExternalID     string           `json:"external_id"`
}

// TransactionsInsertRequest
type TransactionsInsertRequest struct {
	Transactions      []*TransactionInsertItem `json:"transactions"`
	ApplyRules        bool                     `json:"apply_rules,omitempty"`
	SkipDuplicates    bool                     `json:"skip_duplicates,omitempty"`
	CheckForRecurring bool                     `json:"check_for_recurring,omitempty"`
	DebitAsNegative   bool                     `json:"debit_as_negative,omitempty"`
	SkipBalanceUpdate bool                     `json:"skip_balance_update,omitempty"`
}

// TransactionPut
type TransactionInsertItem struct {
	Date        time.Time `json:"date"`
	Amount      string    `json:"amount"`
	CategoryID  *int64    `json:"category_id,omitempty"`
	Payee       string    `json:"payee"`
	Currency    *string   `json:"currency,omitempty"`
	AssetID     int64     `json:"asset_id"`
	RecurringID *int64    `json:"recurring_id,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	Status      *string   `json:"status,omitempty"`
	ExternalID  *string   `json:"external_id,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

// TransactionsUpdateRequest
type TransactionsUpdateRequest struct {
	//Split *TransactionUpdateItem `json:"split"`
	Transaction       *TransactionUpdateItem `json:"transaction"`
	DebitAsNegative   bool                   `json:"debit_as_negative,omitempty"`
	SkipBalanceUpdate bool                   `json:"skip_balance_update,omitempty"`
}

// TransactionUpdateItem
type TransactionUpdateItem struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	Amount      string    `json:"amount"`
	CategoryID  *int64    `json:"category_id,omitempty"`
	Payee       string    `json:"payee"`
	Currency    *string   `json:"currency,omitempty"`
	AssetID     int64     `json:"asset_id"`
	RecurringID *int64    `json:"recurring_id,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	Status      *string   `json:"status,omitempty"`
	ExternalID  *string   `json:"external_id,omitempty"`
	Tags        *[]string `json:"tags,omitempty"`
}

type TransactionsResponse struct {
	Errors *[]string `json:"error,omitempty"`

	// get/list
	Transaction  *Transaction   `json:"transaction,omitempty"`
	Transactions []*Transaction `json:"transactions,omitempty"`

	// insert
	InsertedIds []int64 `json:"ids"`

	// update
	Updated bool     `json:"updated,omitempty"`
	Split   *[]int64 `json:"split,omitempty"`
}

package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/go-playground/validator/v10"
)

// TransactionsResponse is the response we get from requesting transactions.
type TransactionGetResponse struct {
	Transactions []*Transaction `json:"transactions"`
}

// Transaction is a single LM transaction.
type Transaction struct {
	ID             int64  `json:"id"`
	Date           string `json:"date" validate:"datetime=2006-01-02"`
	Payee          string `json:"payee"`
	Amount         string `json:"amount"`
	Currency       string `json:"currency"`
	Notes          string `json:"notes"`
	CategoryID     int64  `json:"category_id"`
	RecurringID    int64  `json:"recurring_id"`
	AssetID        int64  `json:"asset_id"`
	PlaidAccountID int64  `json:"plaid_account_id"`
	Status         string `json:"status"`
	IsGroup        bool   `json:"is_group"`
	GroupID        int64  `json:"group_id"`
	ParentID       int64  `json:"parent_id"`
	ExternalID     string `json:"external_id"`
}

// GetTransactions gets all transactions filtered by the filters.
func (c *Client) GetTransactions(ctx context.Context, filters *TransactionFilters) ([]*Transaction, error) {
	validate := validator.New()
	options := map[string]string{}
	if filters != nil {
		if err := validate.Struct(filters); err != nil {
			return nil, err
		}

		maps, err := filters.ToMap()
		if err != nil {
			return nil, err
		}
		options = maps
	}

	body, err := c.Get(ctx, "/v1/transactions", options)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	resp := &TransactionGetResponse{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp.Transactions, nil
}

// GetTransaction gets a transaction by id.
func (c *Client) GetTransaction(ctx context.Context, id int64, filters *TransactionFilters) (*Transaction, error) {
	validate := validator.New()
	options := map[string]string{}
	if filters != nil {
		if err := validate.Struct(filters); err != nil {
			return nil, err
		}

		maps, err := filters.ToMap()
		if err != nil {
			return nil, err
		}
		options = maps
	}

	body, err := c.Get(ctx, fmt.Sprintf("/v1/transactions/%d", id), options)
	if err != nil {
		return nil, fmt.Errorf("get transaction %d: %w", id, err)
	}

	resp := &Transaction{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// TransactionsInsertRequest
type TransactionsInsertRequest struct {
	Transactions      []*TransactionInsertItem `json:"transactions"`
	ApplyRules        *bool                    `json:"apply_rules,omitempty"`
	SkipDuplicates    *bool                    `json:"skip_duplicates,omitempty"`
	CheckForRecurring *bool                    `json:"check_for_recurring,omitempty"`
	DebitAsNegative   *bool                    `json:"debit_as_negative,omitempty"`
	SkipBalanceUpdate *bool                    `json:"skip_balance_update,omitempty"`
}

// TransactionPut
type TransactionInsertItem struct {
	Date        string   `json:"date" validate:"datetime=2006-01-02"`
	Amount      string   `json:"amount"`
	CategoryID  int64    `json:"category_id,omitempty"`
	Payee       string   `json:"payee"`
	Currency    string   `json:"currency,omitempty"`
	AssetID     int64    `json:"asset_id,omitempty"`
	RecurringID int64    `json:"recurring_id,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Status      string   `json:"status,omitempty"`
	ExternalID  string   `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type TransactionInsertResponse struct {
	InsertedIds []int64 `json:"ids"`
}

// TransactionsInsert updates a transaction by id.
func (c *Client) TransactionsInsert(ctx context.Context, id int64, transactionsInsert TransactionsInsertRequest) (*TransactionInsertResponse, error) {
	validate := validator.New()
	options := map[string]string{}

	updateBody, _ := json.Marshal(transactionsInsert)
	body, err := c.Put(ctx, fmt.Sprintf("/v1/transactions/%d", id), updateBody, options)
	if err != nil {
		return nil, fmt.Errorf("update transaction %d: %w", id, err)
	}

	resp := &TransactionInsertResponse{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp, nil
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
	Date        string   `json:"date,omitempty" validate:"datetime=2006-01-02"`
	Amount      string   `json:"amount,omitempty"`
	CategoryID  int64    `json:"category_id,omitempty"`
	Payee       string   `json:"payee,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	AssetID     int64    `json:"asset_id,omitempty"`
	RecurringID int64    `json:"recurring_id,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Status      string   `json:"status,omitempty"`
	ExternalID  string   `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

//TODO SPLIT

type TransactionUpdateResponse struct {
	Updated bool    `json:"updated"`
	Split   []int64 `json:"split,omitempty"`
}

// UpdateTransaction updates a transaction by id.
func (c *Client) UpdateTransaction(ctx context.Context, id int64, transactionUpdate *TransactionsUpdateRequest) (*TransactionUpdateResponse, error) {
	validate := validator.New()
	options := map[string]string{}

	updateBody, _ := json.Marshal(transactionUpdate)
	body, err := c.Put(ctx, fmt.Sprintf("/v1/transactions/%d", id), updateBody, options)
	if err != nil {
		return nil, fmt.Errorf("update transaction %d: %w", id, err)
	}

	resp := &TransactionUpdateResponse{}
	if err := json.NewDecoder(body).Decode(resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	if err := validate.Struct(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ParsedAmount turns the currency from lunchmoney into a Go currency.
func (t *Transaction) ParsedAmount() (*money.Money, error) {
	return ParseCurrency(t.Amount, t.Currency)
}

// TransactionFilters are options to pass into the request for transactions.
type TransactionFilters struct {
	TagID           *int64  `json:"tag_id,omitempty"`
	RecurringID     *int64  `json:"recurring_id,omitempty"`
	PlaidAccountID  *int64  `json:"plaid_account_id,omitempty"`
	CategoryID      *int64  `json:"category_id,omitempty"`
	AssetID         *string `json:"asset_id,omitempty"`
	Offset          *int64  `json:"offset,omitempty"`
	Limit           *int64  `json:"limit,omitempty"`
	StartDate       *string `json:"start_date,omitempty" validate:"datetime=2006-01-02"`
	EndDate         *string `json:"end_date,omitempty" validate:"datetime=2006-01-02"`
	DebitAsNegative *string `json:"debit_as_negative,omitempty"`
}

// ToMap converts the filters to a string map to be sent with the request as
// GET parameters.
func (r *TransactionFilters) ToMap() (map[string]string, error) {
	ret := map[string]string{}
	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

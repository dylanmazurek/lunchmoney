package models

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/google/go-querystring/query"
)

const (
	defaultFetchLimit    = 40
	defaultFetchPageSize = 10
)

// TransactionFilter are options to pass into the request for transactions.
type TransactionFilter struct {
	assetID int `url:"asset_id"`
	limit   int `url:"limit"`
	offset  int `url:"offset"`

	page     int
	pageSize int

	TagID           *int       `url:"tag_id,omitempty"`
	RecurringID     *int       `url:"recurring_id,omitempty"`
	PlaidAccountID  *int       `url:"plaid_account_id,omitempty"`
	CategoryID      *int       `url:"category_id,omitempty"`
	StartDate       *time.Time `url:"start_date,omitempty" layout:"2006-01-02"`
	EndDate         *time.Time `url:"end_date,omitempty" layout:"2006-01-02"`
	DebitAsNegative bool       `url:"debit_as_negative"`
}

func (t *TransactionFilter) Set(assetID int, pageSize *int, limit *int) {
	t.assetID = assetID

	t.pageSize = defaultFetchPageSize
	if pageSize != nil {
		t.pageSize = *pageSize
	}

	t.limit = defaultFetchLimit
	if limit != nil {
		t.limit = *limit
	}

	t.page = 1
	t.offset = 0
}

func (t *TransactionFilter) NextPage(currCount int) bool {
	t.page++
	t.offset = t.page * t.pageSize

	reachedLimit := t.ReachedLimit(currCount)
	morePages := currCount >= t.pageSize

	return !reachedLimit && morePages
}

func (t *TransactionFilter) ReachedLimit(idx int) bool {
	return idx >= t.limit
}

func (t *TransactionFilter) CurrentPage() int {
	return t.page
}

func (t TransactionFilter) ToVals() (url.Values, error) {
	vals, err := query.Values(t)
	if err != nil {
		return nil, err
	}

	vals.Set("asset_id", fmt.Sprintf("%d", t.assetID))

	if t.limit > 0 {
		vals.Set("limit", fmt.Sprintf("%d", t.limit))
	}

	if t.offset > 0 {
		vals.Set("offset", fmt.Sprintf("%d", t.offset))
	}

	return vals, nil
}

// Transaction is a single transaction.
type Transaction struct {
	ID          *int64           `json:"id,omitempty"`
	Date        DateTimeOptional `json:"date"`
	Payee       string           `json:"payee"`
	Notes       *string          `json:"notes,omitempty"`
	CategoryID  *int64           `json:"category_id,omitempty"`
	RecurringID *int64           `json:"recurring_id,omitempty"`
	AssetID     int64            `json:"asset_id"`
	Status      *string          `json:"status,omitempty"`
	IsGroup     bool             `json:"is_group,omitempty"`
	GroupID     *int64           `json:"group_id,omitempty"`
	ParentID    *int64           `json:"parent_id,omitempty"`
	ExternalID  *string          `json:"external_id,omitempty"`
	Tags        *[]string        `json:"tags,omitempty"`

	Amount money.Money `json:"-"`
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

func (t *Transaction) CreateHash() (hash *Hash) {
	regex := regexp.MustCompile(`([^a-z0-9]+)`)

	if t.Amount.IsZero() {
		return nil
	}

	payee := strings.ToLower(t.Payee)
	payee = regex.ReplaceAllString(payee, "")

	externalId := ""
	if t.ExternalID != nil {
		externalId = *t.ExternalID
	}

	hashString := fmt.Sprintf("%s|%s|%.2f", externalId, payee, t.Amount.AsMajorUnits())

	newHash := Hash(hashString)

	return &newHash
}

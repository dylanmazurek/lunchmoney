package lunchmoney

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dylanmazurek/lunchmoney/models"
)

// func (c *Client) InitTransactionCache(filter models.TransactionFilter) {
// 	transactions, err := c.PageTransactions(context.Background(), &filter)
// 	if err != nil {
// 		return
// 	}

// }

// GetTransaction gets a transaction by id.
func (c *Client) GetTransaction(ctx context.Context, id int64) (*models.Transaction, error) {
	path := fmt.Sprintf("/v1/transactions/%d", id)

	reqOptions := models.Request{
		Method:      http.MethodGet,
		Path:        path,
		QueryValues: nil,
		ReqBody:     nil,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get transaction %d: %w", id, err)
	}

	if respBody.Item == nil {
		return nil, nil
	}

	item := *respBody.Item
	transaction := item.(models.Transaction)

	return &transaction, nil
}

// func (c *Client) FindTransaction(ctx context.Context, transaction models.Transaction, skipCache bool) (*models.Transaction, error) {
// 	hash := transaction.CreateHash()

// 	if hash == nil {
// 		return nil, nil
// 	}

// 	if !skipCache {
// 		if transaction, ok := c.TransactionCache.Cache[*hash]; ok {
// 			return &transaction, nil
// 		}

// 		return nil, nil
// 	}

// 	return nil, nil
// }

func (c *Client) ListTransactions(ctx context.Context, filter *models.TransactionFilter) ([]models.Transaction, error) {
	response, errors := c.StreamTransactions(ctx, filter)

	records := make([]models.Transaction, 0)
	for record := range response {
		records = append(records, record)
	}

	if err := <-errors; err != nil {
		return nil, err
	}

	return records, nil
}

// Retrieve a single page of ListTransactions records from the API. Request is executed immediately.
func (c *Client) PageTransaction(ctx context.Context, vals url.Values) (*models.Response, error) {
	path := "/v1/transactions"

	reqOptions := models.Request{
		Method:      http.MethodGet,
		Path:        path,
		QueryValues: vals,
		ReqBody:     nil,
	}

	resp, err := Request(ctx, c, reqOptions)

	return resp, err
}

// Streams Transaction records from the API as a channel stream. This operation lazily loads records as efficiently as possible until the limit is reached.
func (c *Client) StreamTransactions(ctx context.Context, filter *models.TransactionFilter) (chan models.Transaction, chan error) {
	if filter == nil {
		filter = &models.TransactionFilter{}
	}

	recordChannel := make(chan models.Transaction, 1)
	errorChannel := make(chan error, 1)

	vals, err := filter.ToVals()
	if err != nil {
		errorChannel <- err
		close(recordChannel)
		close(errorChannel)
	}

	firstPage, err := c.PageTransaction(ctx, vals)
	if err != nil {
		errorChannel <- err
		close(recordChannel)
		close(errorChannel)
	} else {
		go c.streamTransaction(ctx, firstPage, filter, recordChannel, errorChannel)
	}

	return recordChannel, errorChannel
}

// streamTransactions
func (c *Client) streamTransaction(ctx context.Context, firstPage *models.Response, filter *models.TransactionFilter, recordChannel chan models.Transaction, errorChannel chan error) {
	transactions := ([]models.Transaction)(*firstPage.Transactions)

	curRecord := 1
	for transactions != nil {
		for itemIndex := range transactions {
			recordChannel <- transactions[itemIndex]
			curRecord += 1
			if filter.ReachedLimit(curRecord) {
				close(recordChannel)
				close(errorChannel)
				return
			}
		}

		vals, err := filter.ToVals()
		if err != nil {
			errorChannel <- err
			break
		}

		transactionsPage, err := c.PageTransaction(context.Background(), vals)
		if err != nil {
			errorChannel <- err
			break
		}

		transactions = ([]models.Transaction)(*transactionsPage.Transactions)
		morePages := filter.NextPage(len(transactions))

		if !morePages {
			break
		}
	}

	close(recordChannel)
	close(errorChannel)
}

type InsertRequest struct {
	Transactions      []models.Transaction `json:"transactions"`
	ApplyRules        bool                 `json:"apply_rules,omitempty"`
	SkipDuplicates    bool                 `json:"skip_duplicates,omitempty"`
	CheckForRecurring bool                 `json:"check_for_recurring,omitempty"`
	DebitAsNegative   bool                 `json:"debit_as_negative,omitempty"`
	SkipBalanceUpdate bool                 `json:"skip_balance_update,omitempty"`
}

func (c *Client) InsertTransactions(ctx context.Context, transactions []models.Transaction, debitAsNegative bool) (*[]int, error) {
	path := "/v1/transactions"

	insertReqBody := &InsertRequest{
		Transactions:      transactions,
		ApplyRules:        true,
		SkipDuplicates:    true,
		CheckForRecurring: true,
		DebitAsNegative:   debitAsNegative,
		SkipBalanceUpdate: true,
	}

	reqOptions := models.Request{
		Method:  http.MethodPost,
		Path:    path,
		ReqBody: insertReqBody,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("insert transactions: %w", err)
	}

	return respBody.Ids, nil
}

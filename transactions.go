package lunchmoney

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dylanmazurek/lunchmoney/models"
)

// GetTransaction gets a transaction by id.
func (c *Client) GetTransaction(ctx context.Context, id int64) (*models.Transaction, error) {
	path := fmt.Sprintf("/v1/transactions/%d", id)

	reqOptions := models.RequestOptions{
		Method:      http.MethodGet,
		Path:        path,
		QueryValues: nil,
		ReqBody:     nil,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get transaction %d: %w", id, err)
	}

	return respBody.Transaction, nil
}

// GetTransactions gets all transactions filtered by the filters.
func (c *Client) GetTransactions(ctx context.Context, filter *models.TransactionFilter) (transactions *[]models.Transaction, err error) {
	path := "/v1/transactions"

	reqOptions := models.RequestOptions{
		Method:      http.MethodGet,
		Path:        path,
		QueryValues: filter,
		ReqBody:     nil,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	return respBody.Transactions, nil
}

// InsertTransactions inserts one or multiple new transactions.
func (c *Client) InsertTransactions(ctx context.Context, body *models.TransactionsInsertRequest) (transactions *[]models.Transaction, wasInserted *bool, err error) {
	path := "/v1/transactions"

	reqOptions := models.RequestOptions{
		Method:  http.MethodPost,
		Path:    path,
		ReqBody: body,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("insert transactions: %w", err)
	}

	inserted := len(respBody.InsertedIds) > 0
	return respBody.Transactions, &inserted, err
}

// UpdateTransaction updates a transaction by id.
func (c *Client) UpdateTransaction(ctx context.Context, transId int64, reqBody *models.TransactionsUpdateRequest) (transaction *models.Transaction, wasUpdated *bool, err error) {
	path := fmt.Sprintf("/v1/transactions/%d", transId)

	reqOptions := models.RequestOptions{
		Method:  http.MethodPut,
		Path:    path,
		ReqBody: reqBody,
	}

	respBody, err := Request(ctx, c, reqOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("update transaction: %w", err)
	}

	return respBody.Transaction, &respBody.Updated, nil
}

package lunchmoney

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

	resp, err := Request[models.TransactionsResponse](ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get transaction %d: %w", id, err)
	}

	if resp.Errors != nil {
		return nil, errors.New(strings.Join(*resp.Errors, ", "))
	}

	return resp.Transaction, nil
}

// GetTransactions gets all transactions filtered by the filters.
func (c *Client) GetTransactions(ctx context.Context, filter *models.TransactionFilter) ([]*models.Transaction, error) {
	path := "/v1/transactions"

	reqOptions := models.RequestOptions{
		Method:      http.MethodGet,
		Path:        path,
		QueryValues: filter,
		ReqBody:     nil,
	}

	resp, err := Request[models.TransactionsResponse](ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}

	if resp.Errors != nil {
		return nil, errors.New(strings.Join(*resp.Errors, ", "))
	}

	return resp.Transactions, nil
}

// InsertTransactions inserts one or multiple new transactions.
func (c *Client) InsertTransactions(ctx context.Context, body *models.TransactionsInsertRequest) (*models.TransactionsResponse, error) {
	path := "/v1/transactions"

	reqOptions := models.RequestOptions{
		Method:      http.MethodPost,
		Path:        path,
		QueryValues: nil,
		ReqBody:     body,
	}

	resp, err := Request[models.TransactionsResponse](ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("insert transactions: %w", err)
	}

	if resp.Errors != nil {
		return nil, errors.New(strings.Join(*resp.Errors, ", "))
	}

	return resp, nil
}

// UpdateTransaction updates a transaction by id.
func (c *Client) UpdateTransaction(ctx context.Context, transId int64, body *models.TransactionsUpdateRequest) (resp *models.TransactionsResponse, err error) {
	path := fmt.Sprintf("/v1/transactions/%d", transId)

	reqOptions := models.RequestOptions{
		Method:      http.MethodPut,
		Path:        path,
		QueryValues: nil,
		ReqBody:     body,
	}

	resp, err = Request[models.TransactionsResponse](ctx, c, reqOptions)
	if err != nil {
		return nil, fmt.Errorf("update transaction: %w", err)
	}

	if resp.Errors != nil {
		return nil, errors.New(strings.Join(*resp.Errors, ", "))
	}

	return resp, nil
}

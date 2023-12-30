package lunchmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

type ListTransactionFilter struct {
}

func (c *Client) ListTransaction(ctx context.Context) (*[]models.Transaction, error) {
	urlString := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Transactions)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, requestUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	query := &url.Values{
		"start": {"01/01/2023"},
		"end":   {"25/12/2023"},
	}

	var transactions models.TransactionResponse
	err = c.Do(ctx, req, &transactions, *query)

	return &transactions.Transactions, err
}

func (c *Client) InsertTransactions(ctx context.Context, transactions []models.Transaction, debitAsNegative bool) (*[]int64, error) {
	insertReqBody := &models.InsertRequest{
		Transactions:      transactions,
		ApplyRules:        true,
		SkipDuplicates:    true,
		CheckForRecurring: true,
		DebitAsNegative:   debitAsNegative,
		SkipBalanceUpdate: true,
	}

	insertJson, err := json.Marshal(&insertReqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, constants.Path.Transactions, bytes.NewReader(insertJson))
	if err != nil {
		return nil, err
	}

	var transactionInsertResponse models.InsertResponse
	err = c.Do(ctx, req, &transactionInsertResponse, nil)

	if transactionInsertResponse.Error != nil {
		err = fmt.Errorf("%s", *transactionInsertResponse.Error)
	}

	return &transactionInsertResponse.Ids, err
}

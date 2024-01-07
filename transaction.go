package lunchmoney

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/util/constants"
)

type ListTransactionFilter struct {
}

func (c *Client) ListTransaction() (*[]models.Transaction, error) {
	urlString := fmt.Sprintf("%s/%s", constants.Config.APIBaseURL, constants.Path.Transactions)
	requestUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	params := &url.Values{
		"start": {"01/01/2023"},
		"end":   {"25/12/2023"},
	}

	req, err := c.NewRequest(http.MethodGet, requestUrl.String(), nil, params)
	if err != nil {
		return nil, err
	}

	var transactions models.TransactionResponse
	err = c.Do(req, &transactions)

	return &transactions.Transactions, err
}

func (c *Client) InsertTransactions(transactions []models.Transaction, debitAsNegative bool) (*[]int64, error) {
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

	req, err := c.NewRequest(http.MethodPost, constants.Path.Transactions, bytes.NewReader(insertJson), nil)
	if err != nil {
		return nil, err
	}

	var transactionInsertResponse models.InsertResponse
	err = c.Do(req, &transactionInsertResponse)

	if len(transactionInsertResponse.Errors) > 0 {
		for _, oneErr := range transactionInsertResponse.Errors {
			newError := errors.New(oneErr)
			err = errors.Join(err, newError)
		}
	}

	return &transactionInsertResponse.Ids, err
}

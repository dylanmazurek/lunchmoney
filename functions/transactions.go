package functions

import (
	"encoding/json"
	"fmt"

	"github.com/Rhymond/go-money"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/shared"
	"github.com/dylanmazurek/lunchmoney/util/time"
	"github.com/rs/zerolog/log"
)

var TransactionStatus = map[string]string{
	"POSTED":    "cleared",
	"PENDING":   "uncleared",
	"cleared":   "cleared",
	"uncleared": "uncleared",
}

func TransactionHandler(lma *lunchmoney.Client, transaction *shared.Transaction) {
	amount := money.NewFromFloat(transaction.Amount, transaction.Currency)

	assetId := json.Number(fmt.Sprintf("%d", transaction.AssetID))
	status := TransactionStatus[transaction.Status]

	lmTransaction := models.Transaction{
		ExternalID: &transaction.ExternalTransactionID,
		Date:       time.Date{Date: transaction.Datetime},
		Payee:      transaction.Description,
		Amount:     *amount,
		AssetID:    assetId,
		Status:     status,
	}

	insertedTransactions, err := lma.InsertTransactions([]models.Transaction{lmTransaction}, true)
	if err != nil {
		log.Error().Err(err).Msg("unable to insert transaction")
	}

	if insertedTransactions != nil {
		log.Info().Msgf("inserted %d transactions", len(*insertedTransactions))
	}
}

package handlers

import (
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
	amountFloat, err := transaction.Amount.Float64()
	if err != nil {
		log.Error().
			Err(err).
			Str("ext-asset-id", transaction.ExternalAssetID).
			Str("ext-transaction-id", transaction.ExternalTransactionID).
			Msg("unable to upsert asset")
	}

	if transaction.AssetID == nil {
		log.Error().
			Str("ext-asset-id", transaction.ExternalAssetID).
			Str("ext-transaction-id", transaction.ExternalTransactionID).
			Msg("unable to upsert transaction, asset id not set")
	}

	amount := money.NewFromFloat(amountFloat, transaction.Currency)

	status := TransactionStatus[transaction.Status]

	lmTransaction := models.Transaction{
		ExternalID: &transaction.ExternalTransactionID,
		Date:       time.Date{Date: transaction.Datetime},
		Payee:      transaction.Description,
		Amount:     *amount,
		AssetID:    *transaction.AssetID,
		Status:     status,
	}

	insertedTransactions, err := lma.InsertTransactions([]models.Transaction{lmTransaction}, true)
	if err != nil {
		log.Error().
			Err(err).
			Str("ext-asset-id", transaction.ExternalAssetID).
			Str("ext-transaction-id", transaction.ExternalTransactionID).
			Msg("unable to upsert transaction")
	}

	if insertedTransactions != nil {
		log.Info().
			Str("ext-asset-id", transaction.ExternalAssetID).
			Str("ext-transaction-id", transaction.ExternalTransactionID).
			Msg("upserted transaction")
	}
}

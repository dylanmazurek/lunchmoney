package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/markkurossi/tabulate"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	client, err := lunchmoney.New(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.InitClient(ctx)
	if err != nil {
		log.Error().Err(err)
		return
	}

	// assetList, err := client.ListAsset(ctx)
	// if err != nil {
	// 	log.Err(err).Msg("failed to list assets")
	// 	return
	// }

	// printAssets(assetList.Assets)

	transactionList, err := client.ListTransaction(ctx)
	if err != nil {
		log.Err(err).Msg("failed to list assets")
		return
	}

	printTransactions(transactionList.Transactions)
}

func printAssets(assets []models.Asset) {
	tab := tabulate.New(tabulate.Unicode)

	tab.Header("ID")
	tab.Header("TYPE")
	tab.Header("NAME")

	tab.Header("BALANCE")
	tab.Header("LAST UPDATED")
	tab.Header("INSTITUTION")

	for _, asset := range assets {
		row := tab.Row()

		row.Column(fmt.Sprintf("%d", asset.ID))
		row.Column(*asset.TypeName)
		row.Column(*asset.Name)
		row.Column(asset.Balance.Display())
		row.Column(asset.BalanceAsOf.Format("2006-01-02"))
		row.Column(asset.InstitutionName)
	}

	tab.Print(os.Stdout)
}

func printTransactions(transactions []models.Transaction) {
	tab := tabulate.New(tabulate.Unicode)

	tab.Header("ID")
	tab.Header("DATE")
	tab.Header("AMOUNT")
	tab.Header("PAYEE")
	tab.Header("ASSET")

	for _, transaction := range transactions {
		row := tab.Row()

		row.Column(fmt.Sprintf("%d", transaction.ID))
		row.Column(transaction.Date.Format("2006-01-02"))
		row.Column(transaction.Amount.Display())
		row.Column(*transaction.Payee)

		val, _ := transaction.AssetID.Float64()
		row.Column(fmt.Sprintf("%.2f", val))
	}

	tab.Print(os.Stdout)
}

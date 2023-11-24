package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
)

func main() {
	ctx := context.Background()

	apiKey := flag.String("key", "", "Lunchmoney API key")
	assetId := flag.Int("asset", 0, "Asset ID for fetching")
	dateFrom := flag.String("from", "2023-09-01", "")
	duration := flag.Duration("duration", 30*24*time.Hour, "")
	flag.Parse()

	if *apiKey == "" {
		fmt.Println("Lunchmoney API key is required")
		return
	}

	client, err := lunchmoney.NewClient(ctx, *apiKey)
	if err != nil {
		fmt.Println(err)
		return
	}

	listAssets(ctx, client)
	listTransactions(ctx, client, *assetId, *dateFrom, *duration)
}

func listAssets(ctx context.Context, client *lunchmoney.Client) {
	assets, err := client.ListAssets(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, asset := range *assets {
		fmt.Printf("asset: [%d] %s - %s\n", asset.ID, asset.Name, asset.InstitutionName)
	}

	fmt.Printf("asset total: %d\n", len(*assets))
}

func listTransactions(ctx context.Context, client *lunchmoney.Client, assetId int, fromString string, duration time.Duration) {
	dateFrom, _ := time.Parse("2006-01-02", fromString)
	dateTo := dateFrom.Add(duration)

	fmt.Printf("transaction between %s - %s\n", dateFrom.Format("2006-01-02"), dateTo.Format("2006-01-02"))

	filter := &models.TransactionFilter{
		StartDate: &dateFrom,
		EndDate:   &dateTo,
	}

	limit := 1000
	filter.Set(assetId, nil, &limit)

	transactions, err := client.ListTransactions(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, transaction := range transactions {
		fmt.Printf("transaction: [%d] %s\n", transaction.ID, transaction.Payee)
	}

	fmt.Printf("trans total: %d\n", len(transactions))
}

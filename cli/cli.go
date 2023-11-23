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

	//listAssets(ctx, client)
	listTransactions(ctx, client)
}

func listAssets(ctx context.Context, client *lunchmoney.Client) {
	assets, err := client.ListAssets(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, asset := range *assets {
		fmt.Printf("asset: %s\n", asset.Name)
	}
}

func listTransactions(ctx context.Context, client *lunchmoney.Client) {
	from, _ := time.Parse("2006-01-02", "2023-09-01")
	to, _ := time.Parse("2006-01-02", "2023-09-30")

	filter := &models.TransactionFilter{
		StartDate: &from,
		EndDate:   &to,
	}

	filter.Set(61147, nil, nil)

	transactions, err := client.ListTransactions(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, transaction := range transactions {
		fmt.Printf("transaction: %s\n", transaction.Payee)
	}
}

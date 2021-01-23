package main

import "fmt"

func main() {
	ctx := context.Background()
	token := os.Getenv("LUNCHMONEY_TOKEN")
	client, _ := lunchmoney.NewClient(token)
	ts, err := client.GetRecurringExpenses(ctx, nil)
	if err != nil {
		log.Panicf("err: %+v", err)
	}

	for _, t := range ts {
		log.Printf("%+v", t)
	}
}

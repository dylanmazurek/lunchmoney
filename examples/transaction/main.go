package main

import (
	"context"
	"log"
	"os"

	"github.com/dylanmazurek/lunchmoney"
)

func main() {
	ctx := context.Background()
	token := os.Getenv("LUNCHMONEY_TOKEN")
	client, _ := lunchmoney.NewClient(token)
	t, err := client.GetTransaction(ctx, 1, nil)
	if err != nil {
		log.Panicf("err: %+v", err)
	}

	log.Printf("%+v", t)
}

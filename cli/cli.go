package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	client, err := lunchmoney.New(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	apiKey, apiKeyOk := os.LookupEnv("API_KEY")

	if apiKeyOk {
		newCredentials := &models.Secrets{
			APIKey: apiKey,
		}
		client.InitClient(ctx, newCredentials)
	}

	err = client.InitClient(ctx, nil)
	if err != nil {
		log.Error().Err(err)
		return
	}

	assets, err := client.ListAsset(ctx)
	if err != nil {
		log.Err(err).Msg("failed to list assets")
		return
	}

	log.Info().Msgf("listed %d assets", len(*assets))

	transactions, err := client.ListTransaction(ctx)
	if err != nil {
		log.Err(err).Msg("failed to list assets")
		return
	}

	log.Info().Msgf("listed %d transactions", len(*transactions))
}

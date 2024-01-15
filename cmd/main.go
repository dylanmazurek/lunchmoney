package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/dylanmazurek/lunchmoney/handlers"
	"github.com/dylanmazurek/lunchmoney/shared"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	buildInfo, _ := debug.ReadBuildInfo()

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Str("go_version", buildInfo.GoVersion).
		Logger()

	log.Logger = logger

	exit := make(chan os.Signal, 1)

	providers := handlers.New(ctx)

	go func() {
		ctx := context.Background()

		assetsSub := providers.RedisClient.Subscribe(ctx, "lunchmoney.assets")
		assetsCh := assetsSub.Channel()
		for msg := range assetsCh {
			var a shared.Asset
			if err := json.Unmarshal([]byte(msg.Payload), &a); err != nil {
				log.Error().Err(err).Msgf("unable to process message from channel: %s", msg.Channel)
			}

			handlers.AssetHandler(providers.LunchmoneyClient, &a)
		}
	}()

	go func() {
		ctx := context.Background()

		transactionsSub := providers.RedisClient.Subscribe(ctx, "lunchmoney.transactions")
		transactionsCh := transactionsSub.Channel()
		for msg := range transactionsCh {
			var t shared.Transaction
			if err := json.Unmarshal([]byte(msg.Payload), &t); err != nil {
				log.Error().Err(err).Msgf("unable to process message from channel: %s", msg.Channel)
			}

			handlers.TransactionHandler(providers.LunchmoneyClient, &t)
		}
	}()

	log.Info().Msg("ready to recieve jobs")

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Info().Msg("closing server")
}

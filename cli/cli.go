package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/functions"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/shared"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New() Server {
	ctx := context.Background()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Error().Err(err).Msg("unable to parse config env")
	}
	services := getServices(ctx, &cfg)

	return Server{
		ctx:      ctx,
		Config:   &cfg,
		Services: *services,
	}
}

func main() {
	server := New()

	exit := make(chan os.Signal, 1)

	assetsSub := server.Services.RedisClient.Subscribe(server.ctx, "lunchmoney.assets")
	assetsCh := assetsSub.Channel()
	for msg := range assetsCh {
		var a shared.Asset
		if err := json.Unmarshal([]byte(msg.Payload), &a); err != nil {
			log.Error().Err(err).Msgf("unable to process message from channel: %s", msg.Channel)
		}

		functions.AssetHandler(server.Services.LunchmoneyClient, &a)
	}

	transactionsSub := server.Services.RedisClient.Subscribe(server.ctx, "lunchmoney.transactions")
	transactionsCh := transactionsSub.Channel()
	for msg := range transactionsCh {
		var t shared.Transaction
		if err := json.Unmarshal([]byte(msg.Payload), &t); err != nil {
			log.Error().Err(err).Msgf("unable to process message from channel: %s", msg.Channel)
		}

		functions.TransactionHandler(server.Services.LunchmoneyClient, &t)
	}

	log.Info().Msg("ready to recieve jobs")

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	log.Info().Msg("closing server")
}

func getServices(ctx context.Context, config *Config) *ServiceProviders {
	lunchmoneyClient, err := lunchmoney.New(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create lunchmoney client")
	}

	var newCredentials *models.Secrets
	if config.APIKey != "" {
		log.Info().Msg("username env var set, setting credentials")
		newCredentials = &models.Secrets{
			APIKey: config.APIKey,
		}
	}

	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(redisOpts)

	err = lunchmoneyClient.InitClient(newCredentials)
	if err != nil {
		log.Info().Msg("api key required")
	}

	return &ServiceProviders{
		LunchmoneyClient: lunchmoneyClient,
		RedisClient:      redisClient,
	}
}

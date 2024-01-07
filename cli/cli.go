package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/functions"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/dylanmazurek/lunchmoney/shared"
	"github.com/dylanmazurek/lunchmoney/util/natsutils"
	"github.com/nats-io/nats.go"
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

	config := server.Config

	exit := make(chan os.Signal, 1)

	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		log.Error().Err(err).Msg("failed to create nats client")
	}

	ec, _ := natsutils.NewEncodedJsonConn(config.NatsUrl)
	defer nc.Close()

	assetsTopic := natsutils.TopicParse("lunchmoney", []string{"assets"})
	natsutils.SubscribeEc(ec, assetsTopic, func(a *shared.Asset) {
		functions.AssetHandler(server.Services.LunchmoneyClient, a)
	})

	transactionsTopic := natsutils.TopicParse("lunchmoney", []string{"transactions"})
	natsutils.SubscribeEc(ec, transactionsTopic, func(t *shared.Transaction) {
		functions.TransactionHandler(server.Services.LunchmoneyClient, t)
	})

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

	err = lunchmoneyClient.InitClient(newCredentials)
	if err != nil {
		log.Info().Msg("api key required")
	}

	return &ServiceProviders{
		LunchmoneyClient: lunchmoneyClient,
	}
}

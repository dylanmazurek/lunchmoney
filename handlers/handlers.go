package handlers

import (
	"context"

	"github.com/caarlos0/env/v10"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/models"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type ServiceProvider struct {
	Config Config

	LunchmoneyClient *lunchmoney.Client
	RedisClient      *redis.Client
}

type Config struct {
	RedisUrl string `env:"REDIS_URL"`

	APIKey *string `env:"API_KEY"`
}

func New(ctx context.Context) *ServiceProvider {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		log.Error().Err(err).Msg("unable to parse config env")
	}

	sp := initServiceProvider(ctx, config)

	return &sp
}

func newLunchmoneyClient(ctx context.Context, config Config) lunchmoney.Client {
	lunchmoneyClient, err := lunchmoney.New(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create wych client")
	}

	var newSecrets *models.Secrets
	if config.APIKey != nil {
		log.Info().Msg("username env var set, setting credentials")
		newSecrets = &models.Secrets{
			APIKey: *config.APIKey,
		}
	}

	err = lunchmoneyClient.InitClient(newSecrets)
	if err != nil {
		log.Error().Err(err).Msg("unable to init client, credentials required")
	}

	log.Info().Msg("wych client initiated")

	return *lunchmoneyClient
}

func newRedisClient(config Config) redis.Client {
	redisOpts, err := redis.ParseURL(config.RedisUrl)
	if err != nil {
		panic(err)
	}

	redisClient := redis.NewClient(redisOpts)

	return *redisClient
}

func initServiceProvider(ctx context.Context, config *Config) ServiceProvider {
	lunchmoneyClient := newLunchmoneyClient(ctx, *config)

	redisClient := newRedisClient(*config)

	return ServiceProvider{
		RedisClient:      &redisClient,
		LunchmoneyClient: &lunchmoneyClient,
	}
}

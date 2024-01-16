package handlers

import (
	"context"

	"github.com/caarlos0/env/v10"
	"github.com/dylanmazurek/lunchmoney"
	"github.com/dylanmazurek/lunchmoney/util/secretstore"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

var (
	SOURCE_NAME = "lunchmoney"
)

type ServiceProvider struct {
	Config Config

	LunchmoneyClient *lunchmoney.Client
	RedisClient      *redis.Client
}

type Config struct {
	RedisUrl string `env:"REDIS_URL"`

	SecretsUrl string `env:"SECRETS_URL"`
	AccessKey  string `env:"ACCESS_KEY"`
}

func New(ctx context.Context) *ServiceProvider {
	log.Info().Msg("loading providers")

	config := Config{}
	if err := env.Parse(&config); err != nil {
		log.Panic().Err(err).Msg("unable to parse config env")
		panic(err)
	}

	sp := initServiceProvider(ctx, config)

	return &sp
}

func newSecretStore(config Config) *secretstore.Client {
	secretStoreClient, err := secretstore.New(config.SecretsUrl, SOURCE_NAME, config.AccessKey)
	if err != nil {
		log.Panic().Err(err).Msg("failed to load secrets")
		panic(err)
	}

	err = secretStoreClient.PingServer()
	if err != nil {
		log.Panic().Err(err).Msg("failed ping secrets server")
		panic(err)
	}

	return secretStoreClient
}

func newLunchmoneyClient(ctx context.Context, secretStore *secretstore.Client) lunchmoney.Client {
	lunchmoneyClient, err := lunchmoney.New(ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to create wych client")
	}

	err = lunchmoneyClient.InitClient(secretStore)
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

func initServiceProvider(ctx context.Context, config Config) ServiceProvider {
	secretStore := newSecretStore(config)
	redisClient := newRedisClient(config)

	lunchmoneyClient := newLunchmoneyClient(ctx, secretStore)

	return ServiceProvider{
		RedisClient:      &redisClient,
		LunchmoneyClient: &lunchmoneyClient,
	}
}

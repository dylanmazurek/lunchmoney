package main

import (
	"context"

	"github.com/dylanmazurek/lunchmoney"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	ctx context.Context

	Config *Config

	Services ServiceProviders
}

type Config struct {
	RedisUrl string `env:"REDIS_URL"`

	APIKey string `env:"API_KEY"`
}

type ServiceProviders struct {
	LunchmoneyClient *lunchmoney.Client

	RedisClient *redis.Client
}

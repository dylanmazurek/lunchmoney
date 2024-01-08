package main

import (
	"context"

	"github.com/dylanmazurek/lunchmoney"
	"github.com/nats-io/nats.go"
)

type Server struct {
	ctx context.Context

	Config *Config

	Services ServiceProviders
}

type Config struct {
	NatsUrl string `env:"NATS_URL"`

	APIKey string `env:"API_KEY"`
}

type ServiceProviders struct {
	LunchmoneyClient *lunchmoney.Client

	NatsClient *nats.EncodedConn
}

package config

import (
	"log"
	"mentor/internal/storage/cache"
	postgres "mentor/internal/storage/db"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	cache.RedisConfig

	AddressServerHTTP string `env:"ADDRESS_SERVER_HTTP" env-required:"true"`
	GRPCPort          int    `env:"GRPC_PORT" env-required:"true"`

	Env string `env:"ENV" env-required:"true"`

	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Fatalf("error reading confug: %s", err.Error())
	}
	return cfg
}

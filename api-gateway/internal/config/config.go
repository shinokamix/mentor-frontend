package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Address string `env:"ADDRESS" env-required:"true"`

	Auth   string `env:"AUTH" env-required:"true"`
	Review string `env:"REVIEW" env-required:"true"`
	Mentor string `env:"MENTOR" env-required:"true"`

	Env string `env:"ENV" env-required:"true"`

	Timeout     time.Duration `env:"TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"60s"`
}

func LoadCOnfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err.Error())
	}

	return cfg
}

package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	KafkaBroker          string `env:"KAFKA_BROKERS" env-required:"true"`
	KafkaTopic           string `env:"KAFKA_TOPIC" env-required:"true"`
	KafkaGroupID         string `env:"KAFKA_GROUP_ID" env-required:"true"`
	MentorServiceAddress string `env:"MENTOR_SERVICE_ADDRESS" env-required:"true"`
}

func LoadConfig() *Config {
	cfg := &Config{}
	err := cleanenv.ReadConfig("./.env", cfg)
	if err != nil {
		log.Fatalf("error reading confug: %s", err.Error())
	}
	return cfg
}

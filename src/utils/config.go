package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

// Config represents global app configuration type.
type Config struct {
	Name       string `env:"NAME" envDefault:"byob-task"`
	DGN        string `env:"DGN" envDefault:"local"`
	Debug      bool   `env:"DEBUG" envDefault:"true"`
	ServerPort int32  `env:"SERVER_PORT" envDefault:"3000"`
	LogLevel   string `env:"LOG_LEVEL" envDefault:"INFO"`

	SentryDSN            string        `env:"SENTRY_DSN" envDefault:"https://sentry.io/example/project"`
	SentryContextTimeout time.Duration `env:"SENTRY_CONTEXT_TIMEOUT" envDefault:"5s"`

	MetricsPort string `env:"METRICS_PORT" envDefault:"9000"`

	YoutubeAPIKey string `env:"YOUTUBE_API_KEY" envDefault:"sample_key"`
}

// NewConfig creates a new Config struct.
func NewConfig() *Config {
	err := godotenv.Load(".env")
	fmt.Print(err)
	if err != nil {
		log.Print("Unable to load .env file. Continuing without loading it...")
	}
	cfg := &Config{}
	if err = env.Parse(cfg); err != nil {
		panic(err)
	}
	return cfg
}

package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"sync"
)

type Config struct {
	AppName          string `envconfig:"APP_NAME"`
	Env              string `envconfig:"APP_ENV"`
	ServerAddress    string `envconfig:"SERVER_ADDRESS"`
	ElasticSearchURL string `envconfig:"ELASTICSEARCH_URL"`
	ElasticSearchKey string `envconfig:"ELASTICSEARCH_KEY"`
}

var (
	config *Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		var cfg Config
		_ = godotenv.Load()
		if err := envconfig.Process("", &cfg); err != nil {
			panic(err)
		}

		config = &cfg
	})
	return config
}

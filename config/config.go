package config

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pelletier/go-toml/v2"
)

type DBConfig struct {
	Path string `toml:"PATH"`
	Name string `toml:"NAME"`
}

type ApplicationConfig struct {
	Port                   int           `toml:"PORT"`
	Development            bool          `toml:"DEVELOPMENT"`
	LogPath                string        `toml:"LOG_PATH"`
	GracefulShutdownPeriod time.Duration `toml:"GRACEFUL_SHUTDOWN_PERIOD"`
	GracefulWarnPeriod     time.Duration `toml:"GRACEFUL_WARN_PERIOD"`
	EnablePprof            bool          `toml:"ENABLE_PPROF"`
}

type ExternalWebsocketConfigOptions struct {
	Url           string        `toml:"URL"`
	OriginUrl     string        `toml:"ORIGIN_URL"`
	RetryTime     time.Duration `toml:"-"`
	RetryWaitTime time.Duration `toml:"-"`
}

type PublisherOptions struct {
	Topic             string
	MaxBufferCapacity int
}

type ConsumerOptions struct {
	Topic          string
	MaxElapsedTime time.Duration
	RetryWaitTime  time.Duration
	MaxWorker      int
}

type ConfigDistributorPubsub struct {
	Consumer  ConsumerOptions
	Publisher PublisherOptions
}

type PubsubConfig struct {
	ConfigDistributor ConfigDistributorPubsub
}

type Config struct {
	Application ApplicationConfig
	DB          struct {
		Clover DBConfig
	}
	External struct {
		Coma struct {
			Websocket ExternalWebsocketConfigOptions
		} `toml:"-"`
	}
	Pubsub PubsubConfig `toml:"-"`

	Auth struct {
		User struct {
			AccessTokenKey       string        `toml:"ACCESS_TOKEN_KEY"`
			RefreshTokenKey      string        `toml:"REFRESH_TOKEN_KEY"`
			AccessTokenDuration  time.Duration `toml:"ACCESS_TOKEN_DURATION"`
			RefreshTokenDuration time.Duration `toml:"REFRESH_TOKEN_DURATION"`
		}
	}
}

var cfg Config
var doOnce sync.Once

func New() Config {
	doOnce.Do(func() {
		configPath := filepath.Join(CFG_PATH, CFG_NAME)
		byt, err := os.ReadFile(configPath)
		if err != nil {
			// set default config
			cfg = defaultConfig()
			data, err := toml.Marshal(cfg)
			if err != nil {
				log.Fatal().Err(err).Msg("cannot marshal config")
				return
			}

			err = os.WriteFile(configPath, data, 0777)
			if err != nil {
				log.Fatal().Err(err).Msg("cannot write config")
				return
			}

			return
		}

		err = toml.Unmarshal(byt, &cfg)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot unmarshaling config")
			return
		}

		cfg.Pubsub = defaultPubsubConfig(PUBSUB_MAX_WORKER, PUBSUB_MAX_BUFFER_CAPACITY)
		cfg.External.Coma.Websocket = defaultExternalComaWSConnection(cfg.Application.Port)
	})

	return cfg
}

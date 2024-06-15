package config

import (
	"crypto/rsa"
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
			PublicKeyLocation    string          `toml:"PUBLIC_KEY_LOCATION"`
			PrivateKeyLocation   string          `toml:"PRIVATE_KEY_LOCATION"`
			PrivateKey           *rsa.PrivateKey `toml:"-"`
			PublicKey            *rsa.PublicKey  `toml:"-"`
			AccessTokenDuration  time.Duration   `toml:"ACCESS_TOKEN_DURATION"`
			RefreshTokenDuration time.Duration   `toml:"REFRESH_TOKEN_DURATION"`
		}
	}
}

var cfg Config
var doOnce sync.Once

func New() Config {
	initConst()

	doOnce.Do(func() {
		// check and create storage dir
		if err := createStorageDirIfNotExist(); err != nil {
			log.Fatal().Err(err).
				Msg("creating data directory")
			return
		}

		configPath := filepath.Join(CONST.CFG_PATH, CONST.CFG_NAME)
		byt, err := os.ReadFile(configPath)
		if err != nil {
			// creating configuration directory
			if err := createCFGDirIfNotExist(); err != nil {
				log.Fatal().Err(err).
					Msg("creating cfg directory")
			}

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

		cfg.Pubsub = defaultPubsubConfig(CONST.PUBSUB_MAX_WORKER, CONST.PUBSUB_MAX_BUFFER_CAPACITY)
		cfg.External.Coma.Websocket = defaultExternalComaWSConnection(cfg.Application.Port)
		cfg.Auth.User.PrivateKey = readRSAPrivateKey()
		cfg.Auth.User.PublicKey = readRSAPublicKey()

		if cfg.Auth.User.PrivateKey == nil || cfg.Auth.User.PublicKey == nil {
			log.Fatal().Msg("PrivateKey or PublicKey is empty")
			return
		}
	})

	return cfg
}

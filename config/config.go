package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/coma/coma/internal/x/rand"
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
	Url       string `toml:"URL"`
	OriginUrl string `toml:"ORIGIN_URL"`
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
		}
	}
	Pubsub PubsubConfig

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

func defaultPubsubConfig(maxWorker int, maxBufferCapacity int) PubsubConfig {
	return PubsubConfig{
		ConfigDistributor: ConfigDistributorPubsub{
			Consumer: ConsumerOptions{
				Topic:         "pubsub:distribute-config",
				RetryWaitTime: 10 * time.Second,
				MaxWorker:     maxWorker,
			},
			Publisher: PublisherOptions{
				Topic:             "pubsub:distribute-config",
				MaxBufferCapacity: maxBufferCapacity,
			},
		},
	}
}

func defaultConfig() Config {
	applicationPort := 5899
	accessTokenKey := rand.RandStr(65)
	refreshTokenKey := rand.RandStr(65)
	return Config{
		Application: ApplicationConfig{
			Port:                   applicationPort,
			Development:            false,
			GracefulShutdownPeriod: 30 * time.Second,
			GracefulWarnPeriod:     30 * time.Second,
			EnablePprof:            false,
		},
		DB: struct{ Clover DBConfig }{
			Clover: DBConfig{
				Path: "./database",
				Name: "localdb",
			},
		},
		External: struct {
			Coma struct {
				Websocket ExternalWebsocketConfigOptions
			}
		}{
			Coma: struct {
				Websocket ExternalWebsocketConfigOptions
			}{
				Websocket: ExternalWebsocketConfigOptions{
					Url:       fmt.Sprintf("127.0.0.1:%d", applicationPort),
					OriginUrl: fmt.Sprintf("127.0.0.1:%d/external-coma-connection", applicationPort),
				},
			},
		},
		Pubsub: defaultPubsubConfig(100000, 1000),
		Auth: struct {
			User struct {
				AccessTokenKey       string        "toml:\"ACCESS_TOKEN_KEY\""
				RefreshTokenKey      string        "toml:\"REFRESH_TOKEN_KEY\""
				AccessTokenDuration  time.Duration "toml:\"ACCESS_TOKEN_DURATION\""
				RefreshTokenDuration time.Duration "toml:\"REFRESH_TOKEN_DURATION\""
			}
		}{
			User: struct {
				AccessTokenKey       string        "toml:\"ACCESS_TOKEN_KEY\""
				RefreshTokenKey      string        "toml:\"REFRESH_TOKEN_KEY\""
				AccessTokenDuration  time.Duration "toml:\"ACCESS_TOKEN_DURATION\""
				RefreshTokenDuration time.Duration "toml:\"REFRESH_TOKEN_DURATION\""
			}{
				AccessTokenKey:       accessTokenKey,
				RefreshTokenKey:      refreshTokenKey,
				AccessTokenDuration:  1 * time.Hour,
				RefreshTokenDuration: 720 * time.Hour,
			},
		},
	}
}

func New(path string) Config {
	doOnce.Do(func() {
		byt, err := os.ReadFile(path)
		if err != nil {
			// set default config
			cfg = defaultConfig()
			data, err := toml.Marshal(cfg)
			if err != nil {
				log.Fatalln("cannot marshal config")
				return
			}

			err = os.WriteFile(path, data, 0644)
			if err != nil {
				log.Fatalln("cannot write config")
				return
			}

			return
		}

		err = toml.Unmarshal(byt, &cfg)
		if err != nil {
			log.Fatalln("cannot unmarshaling config")
			return
		}
	})

	return cfg
}

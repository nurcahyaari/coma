package config

import (
	"fmt"
	"time"

	"github.com/coma/coma/internal/x/rand"
)

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
			} "toml:\"-\""
		}{
			Coma: struct {
				Websocket ExternalWebsocketConfigOptions
			}{
				Websocket: ExternalWebsocketConfigOptions{
					Url:       fmt.Sprintf("ws://127.0.0.1:%d/websocket", applicationPort),
					OriginUrl: fmt.Sprintf("http://127.0.0.1:%d/external-coma-connection", applicationPort),
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

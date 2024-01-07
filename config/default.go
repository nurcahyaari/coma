package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/coma/coma/internal/x/rand"
)

const BASE_PATH = "coma"
const CFG_NAME = "coma.cfg"
const DB_PATH = "database"
const APP_PORT = 5899
const PUBSUB_MAX_WORKER = 1000000
const PUBSUB_MAX_BUFFER_CAPACITY = 1000

func GetBaseWorkingDir(path string) string {
	return filepath.Join(path, BASE_PATH)
}

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

func defaultExternalComaWSConnection(appPort int) ExternalWebsocketConfigOptions {
	return ExternalWebsocketConfigOptions{
		Url:       fmt.Sprintf("ws://127.0.0.1:%d/websocket", appPort),
		OriginUrl: fmt.Sprintf("http://127.0.0.1:%d/external-coma-connection", appPort),
		RetryTime: 60 * time.Second,
	}
}

func defaultConfig(path string) Config {
	accessTokenKey := rand.RandStr(65)
	refreshTokenKey := rand.RandStr(65)
	dbPath := filepath.Join(path, BASE_PATH, DB_PATH)
	return Config{
		Application: ApplicationConfig{
			Port:                   APP_PORT,
			Development:            false,
			GracefulShutdownPeriod: 30 * time.Second,
			GracefulWarnPeriod:     30 * time.Second,
			EnablePprof:            false,
		},
		DB: struct{ Clover DBConfig }{
			Clover: DBConfig{
				Path: dbPath,
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
				Websocket: defaultExternalComaWSConnection(APP_PORT),
			},
		},
		Pubsub: defaultPubsubConfig(PUBSUB_MAX_WORKER, PUBSUB_MAX_BUFFER_CAPACITY),
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

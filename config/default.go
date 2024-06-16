package config

import (
	"crypto/rsa"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/nurcahyaari/coma/internal/x/file"
)

var (
	newConstOnce sync.Once
	CONST        ConstObject
)

type ConstObject struct {
	CFG_NAME                         string
	CFG_PATH                         string
	STORAGE_DIR_PATH                 string
	DB_DIR_NAME                      string
	NIX_STORAGE_PATH                 string
	WIN_STORAGE_PATH                 string
	APP_PORT                         int
	PUBSUB_MAX_WORKER                int
	PUBSUB_MAX_BUFFER_CAPACITY       int
	DEFAULT_RSA_BITSIZE              int
	DEFAULT_RSA_PUBLIC_KEY_LOCATION  string
	DEFAULT_RSA_PRIVATE_KEY_LOCATION string
}

func (c *ConstObject) getStorageDirPath(goos string) {
	// TODO: update later
	switch goos {
	case
		"linux",
		"darwin":
		c.STORAGE_DIR_PATH = c.NIX_STORAGE_PATH
	case "windows":
		c.STORAGE_DIR_PATH = c.WIN_STORAGE_PATH
	}
}

func initConst() {
	newConstOnce.Do(func() {
		if isDevelopment() {
			CONST = NewDevelopmentConstObject()
			return
		}

		CONST = NewConstObject()
	})
}

func NewConstObject() ConstObject {
	cfgName := "coma.cfg"
	cfgPath := "/usr/local/opt/coma"
	dbDirName := "database"
	goos := runtime.GOOS
	co := ConstObject{
		CFG_NAME:                         cfgName,
		CFG_PATH:                         cfgPath,
		DB_DIR_NAME:                      dbDirName,
		NIX_STORAGE_PATH:                 "/var/lib/coma",
		APP_PORT:                         5899,
		PUBSUB_MAX_WORKER:                1000000,
		PUBSUB_MAX_BUFFER_CAPACITY:       1000,
		DEFAULT_RSA_BITSIZE:              2048,
		DEFAULT_RSA_PUBLIC_KEY_LOCATION:  cfgPath + "/coma_public.pem",
		DEFAULT_RSA_PRIVATE_KEY_LOCATION: cfgPath + "/coma_private.pem",
	}

	co.getStorageDirPath(goos)

	return co
}

func NewDevelopmentConstObject() ConstObject {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	wd = filepath.Join(wd, "temporary_storage")
	cfgName := "coma.cfg"
	dbPath := "database"
	goos := runtime.GOOS

	co := ConstObject{
		CFG_NAME:                         cfgName,
		CFG_PATH:                         wd,
		DB_DIR_NAME:                      dbPath,
		NIX_STORAGE_PATH:                 wd,
		WIN_STORAGE_PATH:                 wd,
		APP_PORT:                         5898,
		PUBSUB_MAX_WORKER:                1000000,
		PUBSUB_MAX_BUFFER_CAPACITY:       1000,
		DEFAULT_RSA_BITSIZE:              2048,
		DEFAULT_RSA_PUBLIC_KEY_LOCATION:  wd + "/coma_public.pem",
		DEFAULT_RSA_PRIVATE_KEY_LOCATION: wd + "/coma_private.pem",
	}

	co.getStorageDirPath(goos)

	return co
}

func isDevelopment() bool {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	dir := filepath.Dir(ex)

	return strings.Contains(dir, "go-build")
}

// init storage dir creates all data that will store inside storage location
// containing database and keys
func createStorageDirIfNotExist() error {
	if _, err := os.Stat(CONST.STORAGE_DIR_PATH); !os.IsNotExist(err) {
		return nil
	}

	// create base dir for storage
	if err := file.NewDir(CONST.STORAGE_DIR_PATH); err != nil {
		return err
	}

	// create dir for database
	if err := file.NewDir(filepath.Join(CONST.STORAGE_DIR_PATH, CONST.DB_DIR_NAME)); err != nil {
		return err
	}

	// create rsa
	if err := createDefaultRSAIfNotExist(); err != nil {
		return err
	}

	return nil
}

func createCFGDirIfNotExist() error {
	if _, err := os.Stat(CONST.STORAGE_DIR_PATH); !os.IsNotExist(err) {
		return nil
	}

	// create dir for configuration
	if err := file.NewDir(CONST.CFG_PATH); err != nil {
		return err
	}
	return nil
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

func defaultConfig() Config {
	dbPath := filepath.Join(CONST.STORAGE_DIR_PATH, CONST.DB_DIR_NAME)

	return Config{
		Application: ApplicationConfig{
			Port:                   CONST.APP_PORT,
			Development:            isDevelopment(),
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
				Websocket: defaultExternalComaWSConnection(CONST.APP_PORT),
			},
		},
		Pubsub: defaultPubsubConfig(CONST.PUBSUB_MAX_WORKER, CONST.PUBSUB_MAX_BUFFER_CAPACITY),
		Auth: struct {
			User struct {
				PublicKeyLocation    string          "toml:\"PUBLIC_KEY_LOCATION\""
				PrivateKeyLocation   string          "toml:\"PRIVATE_KEY_LOCATION\""
				PrivateKey           *rsa.PrivateKey "toml:\"-\""
				PublicKey            *rsa.PublicKey  "toml:\"-\""
				AccessTokenDuration  time.Duration   "toml:\"ACCESS_TOKEN_DURATION\""
				RefreshTokenDuration time.Duration   "toml:\"REFRESH_TOKEN_DURATION\""
			}
		}{
			User: struct {
				PublicKeyLocation    string          "toml:\"PUBLIC_KEY_LOCATION\""
				PrivateKeyLocation   string          "toml:\"PRIVATE_KEY_LOCATION\""
				PrivateKey           *rsa.PrivateKey "toml:\"-\""
				PublicKey            *rsa.PublicKey  "toml:\"-\""
				AccessTokenDuration  time.Duration   "toml:\"ACCESS_TOKEN_DURATION\""
				RefreshTokenDuration time.Duration   "toml:\"REFRESH_TOKEN_DURATION\""
			}{
				PrivateKeyLocation:   CONST.DEFAULT_RSA_PUBLIC_KEY_LOCATION,
				PublicKeyLocation:    CONST.DEFAULT_RSA_PRIVATE_KEY_LOCATION,
				PrivateKey:           readRSAPrivateKey(),
				PublicKey:            readRSAPublicKey(),
				AccessTokenDuration:  1 * time.Hour,
				RefreshTokenDuration: 720 * time.Hour,
			},
		},
	}
}

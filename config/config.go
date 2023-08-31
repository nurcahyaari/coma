package config

import (
	"log"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Application struct {
		Port        int  `mapstructure:"PORT"`
		Development bool `mapstructure:"DEVELOPMENT"`
		Log         struct {
			Path string `mapstructure:"PATH"`
		} `mapstructure:"LOG"`
		Key struct {
			Default string `mapstructure:"DEFAULT"`
			Rsa     struct {
				Public  string `mapstructure:"PUBLIC"`
				Private string `mapstructure:"PRIVATE"`
			}
		} `mapstructure:"KEY"`
		Graceful struct {
			ShutdownPeriod time.Duration `mapstructure:"SHUTDOWN_PERIOD"`
			WarnPeriod     time.Duration `mapstructure:"WARN_PERIOD"`
		} `mapstructure:"GRACEFUL"`
		Pprof struct {
			Enable bool `mapstructure:"ENABLE"`
		} `mapstructure:"PPROF"`
	} `mapstructure:"APPLICATION"`

	DB struct {
		Mysql struct {
			Host string `mapstructure:"HOST"`
			Port int    `mapstructure:"PORT"`
			Name string `mapstructure:"NAME"`
			User string `mapstructure:"USER"`
			Pass string `mapstructure:"PASS"`
		} `mapstructure:"MYSQL"`

		Clover struct {
			Path string `mapstructure:"PATH"`
			Name string `mapstructure:"NAME"`
		} `mapstructure:"CLOVER"`
	} `mapstructure:"DB"`

	External struct {
		Coma struct {
			Websocket struct {
				Url       string `mapstructure:"URL"`
				OriginUrl string `mapstructure:"ORIGIN_URL"`
			} `mapstructure:"WEBSOCKET"`
		} `mapstructure:"COMA"`
	} `mapstructure:"EXTERNAL"`

	Pubsub struct {
		Local struct {
			Publisher struct {
				ConfigDistributor struct {
					Topic             string `mapstructure:"TOPIC"`
					MaxBufferCapacity int    `mapstructure:"MAX_BUFFER_CAPACITY"`
				} `mapstructure:"CONFIG_DISTRIBUTOR"`
			} `mapstructure:"PUBLISHER"`
			Consumer struct {
				ConfigDistributor struct {
					Topic          string        `mapstructure:"TOPIC"`
					MaxElapsedTime time.Duration `mapstructure:"MAX_ELAPSED_TIME"`
					RetryWaitTime  time.Duration `mapstructure:"RETRY_WAIT_TIME"`
					MaxWorker      int           `mapstructure:"MAX_WORKER"`
				} `mapstructure:"CONFIG_DISTRIBUTOR"`
			} `mapstructure:"CONSUMER"`
		} `mapstructure:"LOCAL"`
	} `mapstructure:"PUBSUB"`

	Auth struct {
		User struct {
			AccessTokenKey       string        `mapstructure:"ACCESS_TOKEN_KEY"`
			RefreshTokenKey      string        `mapstructure:"REFRESH_TOKEN_KEY"`
			AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
			RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
		} `mapstructure:"USER"`
	} `mapstructure:"AUTH"`
}

var cfg Config
var doOnce sync.Once

func Get() Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read .env file")
	}

	doOnce.Do(func() {
		err := viper.Unmarshal(&cfg)
		if err != nil {
			log.Fatalln("cannot unmarshaling config")
		}
	})

	return cfg
}

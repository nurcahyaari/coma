package database

import (
	"path/filepath"

	"github.com/ostafen/clover"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Path string
	Name string
}

func (c Config) BuildDirPath() string {
	return filepath.Join(c.Path, c.Name)
}

type Clover struct {
	DB *clover.DB
}

func NewClover(cfg Config) *Clover {
	db, err := clover.Open(cfg.BuildDirPath())
	if err != nil {
		log.Err(err).Msgf("Error to loading Clover DB %s", err)
		panic(err)
	}

	return &Clover{
		DB: db,
	}
}

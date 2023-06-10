package repository

import "github.com/coma/coma/infrastructure/database"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryWriter
type RepositoryWriter interface{}

//counterfeiter:generate . RepositoryWrite
type RepositoryWrite struct {
	db *database.Clover
}

func NewRepositoryWriter(db *database.Clover) RepositoryWriter {
	return &RepositoryWrite{
		db: db,
	}
}

package repository

import "github.com/coma/coma/infrastructure/database"

type RepositoryApplicationWriter interface {
}

type RepositoryApplicationWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationWriter(db *database.Clover, name string) RepositoryApplicationWriter {
	return &RepositoryApplicationWrite{
		db:     db,
		dbName: name,
	}
}

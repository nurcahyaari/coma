package repository

import "github.com/coma/coma/infrastructure/database"

type RepositoryApplicationReader interface {
}

type RepositoryApplicationRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationReader(db *database.Clover, name string) RepositoryApplicationReader {
	return &RepositoryApplicationRead{
		db:     db,
		dbName: name,
	}
}

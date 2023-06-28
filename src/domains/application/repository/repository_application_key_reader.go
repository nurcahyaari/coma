package repository

import "github.com/coma/coma/infrastructure/database"

type RepositoryApplicationKeyReader interface {
}

type RepositoryApplicationKeyRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationKeyReader(db *database.Clover, name string) RepositoryApplicationKeyReader {
	return &RepositoryApplicationKeyRead{
		db:     db,
		dbName: name,
	}
}

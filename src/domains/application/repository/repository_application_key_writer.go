package repository

import "github.com/coma/coma/infrastructure/database"

type RepositoryApplicationKeyWriter interface {
}

type RepositoryApplicationKeyWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationKeyWriter(db *database.Clover, name string) RepositoryApplicationKeyWriter {
	return &RepositoryApplicationKeyWrite{
		db:     db,
		dbName: name,
	}
}

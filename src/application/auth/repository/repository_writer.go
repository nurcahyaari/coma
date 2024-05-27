package repository

import (
	"github.com/nurcahyaari/coma/infrastructure/database"
	"github.com/nurcahyaari/coma/src/domain/repository"
)

type RepositoryWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryWriter(db *database.Clover, name string) repository.RepositoryAuthWriter {
	db.DB.CreateCollection(name)
	return &RepositoryWrite{
		db:     db,
		dbName: name,
	}
}

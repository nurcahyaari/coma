package repository

import (
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/repository"
)

//counterfeiter:generate . RepositoryWrite
type RepositoryWrite struct {
	db *database.Clover
}

func NewRepositoryWriter(db *database.Clover) repository.RepositoryAuthWriter {
	return &RepositoryWrite{
		db: db,
	}
}

package repository

import (
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/repository"
)

//counterfeiter:generate . Repository
type Repository struct {
	db *database.Clover
}

func New(db *database.Clover) repository.AuthRepositorier {
	db.DB.CreateCollection("apikey")
	return &Repository{
		db: db,
	}
}

func (r Repository) NewRepositoryReader() repository.RepositoryAuthReader {
	return NewRepositoryReader(r.db)
}

func (r Repository) NewRepositoryWriter() repository.RepositoryAuthWriter {
	return NewRepositoryWriter(r.db)
}

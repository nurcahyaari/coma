package repository

import (
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/repository"
)

type Repository struct {
	dbName string
	db     *database.Clover
}

func NewRepository(db *database.Clover) *Repository {
	dbName := "application"
	db.DB.CreateCollection(dbName)
	return &Repository{
		db:     db,
		dbName: dbName,
	}
}

func (r *Repository) NewRepositoryUserReader() repository.RepositoryUserReader {
	return NewRepositoryUserReader(r.db)
}

func (r *Repository) NewRepositoryUserWriter() repository.RepositoryUserWriter {
	return NewRepositoryUserWriter(r.db)
}

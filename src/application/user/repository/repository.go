package repository

import (
	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/repository"
)

type Repository struct {
	name string
	db   *database.Clover
}

func New(db *database.Clover) *Repository {
	dbName := "user"
	return &Repository{
		db:   db,
		name: dbName,
	}
}

func (r *Repository) NewRepositoryUserReader() repository.RepositoryUserReader {
	return NewRepositoryUserReader(r.name, r.db)
}

func (r *Repository) NewRepositoryUserWriter() repository.RepositoryUserWriter {
	return NewRepositoryUserWriter(r.name, r.db)
}

func (r *Repository) NewRepositoryUserApplicationScopeReader() repository.RepositoryUserApplicationScopeReader {
	return NewRepositoryUserApplicationScopeRead(r.name, r.db)
}

func (r *Repository) NewRepositoryUserApplicationScopeWriter() repository.RepositoryUserApplicationScopeWriter {
	return NewRepositoryUserApplicationScopeWrite(r.name, r.db)
}

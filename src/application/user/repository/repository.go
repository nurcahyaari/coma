package repository

import (
	"fmt"

	"github.com/nurcahyaari/coma/infrastructure/database"
	"github.com/nurcahyaari/coma/src/domain/repository"
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
	return NewRepositoryUserApplicationScopeRead(fmt.Sprintf("%s_application_scope", r.name), r.db)
}

func (r *Repository) NewRepositoryUserApplicationScopeWriter() repository.RepositoryUserApplicationScopeWriter {
	return NewRepositoryUserApplicationScopeWrite(fmt.Sprintf("%s_application_scope", r.name), r.db)
}

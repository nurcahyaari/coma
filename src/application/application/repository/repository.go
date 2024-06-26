package repository

import (
	"fmt"

	"github.com/nurcahyaari/coma/infrastructure/database"
	"github.com/nurcahyaari/coma/src/domain/repository"
)

type Repository struct {
	dbName string
	db     *database.Clover
}

func New(db *database.Clover) *Repository {
	dbName := "application"
	return &Repository{
		db:     db,
		dbName: dbName,
	}
}

func (r Repository) NewRepositoryApplicationReader() repository.RepositoryApplicationReader {
	return NewRepositoryApplicationReader(r.db, r.dbName)
}

func (r Repository) NewRepositoryApplicationWriter() repository.RepositoryApplicationWriter {
	return NewRepositoryApplicationWriter(r.db, r.dbName)
}

func (r Repository) NewRepositoryApplicationKeyReader() repository.RepositoryApplicationKeyReader {
	return NewRepositoryApplicationKeyReader(r.db, fmt.Sprintf("%s_key", r.dbName))
}

func (r Repository) NewRepositoryApplicationKeyWriter() repository.RepositoryApplicationKeyWriter {
	return NewRepositoryApplicationKeyWriter(r.db, fmt.Sprintf("%s_key", r.dbName))
}

func (r Repository) NewRepositoryApplicationConfigurationReader() repository.RepositoryApplicationConfigurationReader {
	return NewApplicationConfigurationRepositoryReader(r.db, fmt.Sprintf("%s_configuration", r.dbName))
}

func (r Repository) NewRepositoryApplicationConfigurationWriter() repository.RepositoryApplicationConfigurationWriter {
	return NewApplicationConfigurationRepositoryWriter(r.db, fmt.Sprintf("%s_configuration", r.dbName))
}

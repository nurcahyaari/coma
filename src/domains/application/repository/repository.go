package repository

import (
	"fmt"

	"github.com/coma/coma/infrastructure/database"
)

type Repository struct {
	dbName string
	db     *database.Clover
}

func New(db *database.Clover) *Repository {
	dbName := "application"
	db.DB.CreateCollection(dbName)
	return &Repository{
		db:     db,
		dbName: dbName,
	}
}

func (r Repository) NewRepositoryApplicationReader() RepositoryApplicationReader {
	return NewRepositoryApplicationReader(r.db, r.dbName)
}

func (r Repository) NewRepositoryApplicationWriter() RepositoryApplicationWriter {
	return NewRepositoryApplicationWriter(r.db, r.dbName)
}

func (r Repository) NewRepositoryApplicationStageReader() RepositoryApplicationStageReader {
	return NewRepositoryApplicationStageReader(r.db, fmt.Sprintf("%s_stage", r.dbName))
}

func (r Repository) NewRepositoryApplicationStageWriter() RepositoryApplicationStageWriter {
	return NewRepositoryApplicationStageWriter(r.db, fmt.Sprintf("%s_stage", r.dbName))
}

func (r Repository) NewRepositoryApplicationKeyReader() RepositoryApplicationKeyReader {
	return NewRepositoryApplicationKeyReader(r.db, fmt.Sprintf("%s_key", r.dbName))
}

func (r Repository) NewRepositoryApplicationKeyWriter() RepositoryApplicationKeyWriter {
	return NewRepositoryApplicationKeyWriter(r.db, fmt.Sprintf("%s_key", r.dbName))
}

func (r Repository) NewRepositoryApplicationConfigurationReader() RepositoryApplicationConfigurationReader {
	return NewApplicationConfigurationRepositoryReader(r.db, fmt.Sprintf("%s_configuration", r.dbName))
}

func (r Repository) NewRepositoryApplicationConfigurationWriter() RepositoryApplicationConfigurationWriter {
	return NewApplicationConfigurationRepositoryWriter(r.db, fmt.Sprintf("%s_configuration", r.dbName))
}

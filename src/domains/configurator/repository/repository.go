package repository

import (
	"github.com/coma/coma/infrastructure/database"
)

type Repository struct {
	dbName string
	db     *database.Clover
}

func New(db *database.Clover) *Repository {
	dbName := "configuration"
	db.DB.CreateCollection(dbName)
	return &Repository{
		db:     db,
		dbName: dbName,
	}
}

func (r Repository) NewRepositoryReader() RepositoryReader {
	return NewRepositoryReader(r.db, r.dbName)
}

func (r Repository) NewRepositoryWriter() RepositoryWriter {
	return NewRepositoryWriter(r.db, r.dbName)
}

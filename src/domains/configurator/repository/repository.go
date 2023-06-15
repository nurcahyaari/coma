package repository

import "github.com/coma/coma/infrastructure/database"

type Repository struct {
	db *database.Clover
}

func New(db *database.Clover) *Repository {
	db.DB.CreateCollection("configuration")
	return &Repository{
		db: db,
	}
}

func (r Repository) NewRepositoryReader() RepositoryReader {
	return NewRepositoryReader(r.db)
}

func (r Repository) NewRepositoryWriter() RepositoryWriter {
	return NewRepositoryWriter(r.db)
}

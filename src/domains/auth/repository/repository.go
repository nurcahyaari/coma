package repository

import "github.com/coma/coma/infrastructure/database"

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . Repositorier
type Repositorier interface {
	NewRepositoryReader() RepositoryReader
	NewRepositoryWriter() RepositoryWriter
}

//counterfeiter:generate . Repository
type Repository struct {
	db *database.Clover
}

func New(db *database.Clover) *Repository {
	db.DB.CreateCollection("apikey")
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

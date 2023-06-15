package repository

import "github.com/coma/coma/infrastructure/database"

type RepositoryWriter interface {
}

type RepositoryWrite struct {
	db *database.Clover
}

func NewRepositoryWriter(db *database.Clover) RepositoryWriter {
	db.DB.CreateCollection("configuration")
	return &RepositoryWrite{
		db: db,
	}
}

package repository

import "github.com/coma/coma/infrastructure/database"

type RepositoryReader interface {
}

type RepositoryRead struct {
	db *database.Clover
}

func NewRepositoryReader(db *database.Clover) RepositoryReader {
	return &RepositoryRead{
		db: db,
	}
}

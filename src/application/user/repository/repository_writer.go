package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
)

type RepositoryUserWrite struct {
	db *database.Clover
}

func NewRepositoryUserWriter(db *database.Clover) repository.RepositoryUserWriter {
	return &RepositoryUserWrite{
		db: db,
	}
}

func (r *RepositoryUserWrite) SaveUser(ctx context.Context, user entity.User) error {
	return nil
}

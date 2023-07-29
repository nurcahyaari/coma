package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
)

type RepositoryUserRead struct {
	db *database.Clover
}

func NewRepositoryUserReader(db *database.Clover) repository.RepositoryUserReader {
	return &RepositoryUserRead{
		db: db,
	}
}

func (r *RepositoryUserRead) FindUser(ctx context.Context, filter entity.FilterUser) (entity.User, error)

func (r *RepositoryUserRead) FindUsers(ctx context.Context, filter entity.FilterUser) (entity.Users, error)

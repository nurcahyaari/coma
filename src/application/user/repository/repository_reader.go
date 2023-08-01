package repository

import (
	"context"
	"errors"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
)

type RepositoryUserRead struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserReader(name string, db *database.Clover) repository.RepositoryUserReader {
	return &RepositoryUserRead{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserRead) FindUser(ctx context.Context, filter entity.FilterUser) (entity.User, error) {
	var user entity.User

	users, err := r.FindUsers(ctx, filter)
	if err != nil {
		return user, err
	}
	if len(users) == 0 {
		return user, errors.New("err: user doesn't found")
	}

	user = users[0]

	return user, nil
}

func (r *RepositoryUserRead) FindUsers(ctx context.Context, filter entity.FilterUser) (entity.Users, error) {
	var users entity.Users

	docs, err := r.db.DB.Query(r.name).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return users, err
	}

	for _, doc := range docs {
		user := entity.User{}
		if err := doc.Unmarshal(&user); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
)

type RepositoryUserRead struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserReader(name string, db *database.Clover) repository.RepositoryUserReader {
	db.DB.CreateCollection(name)
	return &RepositoryUserRead{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserRead) FindUser(ctx context.Context, filter entity.FilterUser) (entity.User, error) {
	var user entity.User

	if filter.Filter() == nil {
		return user, nil
	}

	users, err := r.FindUsers(ctx, filter)
	if err != nil {
		internalerrors.StackTrace(err)
		return user, err
	}
	if len(users) == 0 {
		return user, nil
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
		internalerrors.StackTrace(err)
		return users, err
	}

	for _, doc := range docs {
		user := entity.User{}
		if err := doc.Unmarshal(&user); err != nil {
			internalerrors.StackTrace(err)
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

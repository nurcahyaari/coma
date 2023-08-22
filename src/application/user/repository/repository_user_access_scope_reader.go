package repository

import (
	"context"
	"errors"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
)

type RepositoryUserAccessScopeRead struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserAccessScopeRead(name string, db *database.Clover) repository.RepositoryUserAccessScopeReader {
	db.DB.CreateCollection(name)
	return &RepositoryUserAccessScopeRead{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserAccessScopeRead) FindUserAccessScope(ctx context.Context, filter entity.FilterUserAccessScope) (entity.UserAccessScope, error) {
	var userAccessScope entity.UserAccessScope

	if filter.Filter() == nil {
		return userAccessScope, nil
	}

	userAccessesScope, err := r.FindUserAccessesScope(ctx, filter)
	if err != nil {
		return userAccessScope, err
	}
	if len(userAccessesScope) == 0 {
		return userAccessScope, errors.New("err: user access doesn't found")
	}

	userAccessScope = userAccessesScope[0]

	return userAccessScope, nil
}

func (r *RepositoryUserAccessScopeRead) FindUserAccessesScope(ctx context.Context, filter entity.FilterUserAccessScope) (entity.UserAccessesScope, error) {
	var userAccessesScope entity.UserAccessesScope

	docs, err := r.db.DB.Query(r.name).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return userAccessesScope, err
	}

	for _, doc := range docs {
		userAccessScope := entity.UserAccessScope{}
		if err := doc.Unmarshal(&userAccessScope); err != nil {
			return userAccessesScope, err
		}

		userAccessesScope = append(userAccessesScope, userAccessScope)
	}

	return userAccessesScope, nil
}

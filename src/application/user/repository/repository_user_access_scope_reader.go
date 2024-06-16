package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
)

type RepositoryUserApplicationScopeRead struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserApplicationScopeRead(name string, db *database.Clover) repository.RepositoryUserApplicationScopeReader {
	db.DB.CreateCollection(name)
	return &RepositoryUserApplicationScopeRead{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserApplicationScopeRead) FindUserApplicationScope(ctx context.Context, filter entity.FilterUserApplicationScope) (entity.UserApplicationScope, bool, error) {
	var userApplicationScope entity.UserApplicationScope

	if filter.Filter() == nil {
		return userApplicationScope, false, nil
	}

	userApplicationsScope, err := r.FindUserApplicationsScope(ctx, filter)
	if err != nil {
		internalerrors.StackTrace(err)
		return userApplicationScope, false, err
	}
	if len(userApplicationsScope) == 0 {
		return userApplicationScope, false, nil
	}

	userApplicationScope = userApplicationsScope[0]

	return userApplicationScope, true, nil
}

func (r *RepositoryUserApplicationScopeRead) FindUserApplicationsScope(ctx context.Context, filter entity.FilterUserApplicationScope) (entity.UserApplicationsScope, error) {
	var userApplicationsScope entity.UserApplicationsScope

	docs, err := r.db.DB.Query(r.name).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		internalerrors.StackTrace(err)
		return userApplicationsScope, err
	}

	for _, doc := range docs {
		userApplicationScope := entity.UserApplicationScope{}
		if err := doc.Unmarshal(&userApplicationScope); err != nil {
			internalerrors.StackTrace(err)
			return userApplicationsScope, err
		}

		userApplicationsScope = append(userApplicationsScope, userApplicationScope)
	}

	return userApplicationsScope, nil
}

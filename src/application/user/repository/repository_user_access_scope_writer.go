package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryUserAccessScopeWrite struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserAccessScopeWrite(name string, db *database.Clover) repository.RepositoryUserAccessScopeWriter {
	db.DB.CreateCollection(name)
	return &RepositoryUserAccessScopeWrite{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserAccessScopeWrite) SaveUserAccessScope(ctx context.Context, userAccessScope entity.UserAccessScope) error {
	dataMap, err := userAccessScope.MapStringInterface()
	if err != nil {
		return err
	}

	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	_, err = r.db.DB.InsertOne(r.name, doc)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryUserAccessScopeWrite) UpdateUserAccessScope(ctx context.Context, userAccessScope entity.UserAccessScope) error {
	dataMap, err := userAccessScope.MapStringInterface()
	if err != nil {
		return err
	}

	err = r.db.DB.Query(r.name).
		UpdateById(userAccessScope.Id, dataMap)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryUserAccessScopeWrite) DeleteUserAccess(ctx context.Context, filter entity.FilterUserAccessScope) error {
	err := r.db.DB.Query(r.name).
		Where(filter.Filter()).
		Delete()
	if err != nil {
		return err
	}

	return nil
}

package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryUserApplicationScopeWrite struct {
	name string
	db   *database.Clover
}

func NewRepositoryUserApplicationScopeWrite(name string, db *database.Clover) repository.RepositoryUserApplicationScopeWriter {
	db.DB.CreateCollection(name)
	return &RepositoryUserApplicationScopeWrite{
		name: name,
		db:   db,
	}
}

func (r *RepositoryUserApplicationScopeWrite) SaveUserApplicationScope(ctx context.Context, userApplicationScope entity.UserApplicationScope) error {
	dataMap, err := userApplicationScope.MapStringInterface()
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

func (r *RepositoryUserApplicationScopeWrite) UpdateUserApplicationScope(ctx context.Context, userApplicationScope entity.UserApplicationScope) error {
	dataMap, err := userApplicationScope.MapStringInterface()
	if err != nil {
		return err
	}

	err = r.db.DB.Query(r.name).
		UpdateById(userApplicationScope.Id, dataMap)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryUserApplicationScopeWrite) RevokeUserApplicationScope(ctx context.Context, filter entity.FilterUserApplicationScope) error {
	err := r.db.DB.Query(r.name).
		Where(filter.Filter()).
		Delete()
	if err != nil {
		return err
	}

	return nil
}

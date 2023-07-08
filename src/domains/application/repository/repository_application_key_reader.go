package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
)

type RepositoryApplicationKeyReader interface {
	FindApplicationKey(ctx context.Context, filter model.FilterApplicationKey) (model.ApplicationKey, error)
}

type RepositoryApplicationKeyRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationKeyReader(db *database.Clover, name string) RepositoryApplicationKeyReader {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationKeyRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationKeyRead) FindApplicationKey(ctx context.Context, filter model.FilterApplicationKey) (model.ApplicationKey, error) {
	var applicationKey model.ApplicationKey

	doc, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindFirst()
	if err != nil {
		return applicationKey, err
	}
	if doc == nil {
		return applicationKey, nil
	}

	err = doc.Unmarshal(&applicationKey)
	if err != nil {
		return applicationKey, err
	}

	return applicationKey, nil
}

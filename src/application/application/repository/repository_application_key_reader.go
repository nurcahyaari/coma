package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
)

type RepositoryApplicationKeyRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationKeyReader(db *database.Clover, name string) repository.RepositoryApplicationKeyReader {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationKeyRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationKeyRead) FindApplicationKey(ctx context.Context, filter entity.FilterApplicationKey) (entity.ApplicationKey, error) {
	var applicationKey entity.ApplicationKey

	doc, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindFirst()
	if err != nil {
		internalerrors.StackTrace(err)
		return applicationKey, err
	}
	if doc == nil {
		return applicationKey, nil
	}

	err = doc.Unmarshal(&applicationKey)
	if err != nil {
		internalerrors.StackTrace(err)
		return applicationKey, err
	}

	return applicationKey, nil
}

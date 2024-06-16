package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryApplicationConfigurationWrite struct {
	dbName string
	db     *database.Clover
}

func NewApplicationConfigurationRepositoryWriter(db *database.Clover, name string) repository.RepositoryApplicationConfigurationWriter {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationConfigurationWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationConfigurationWrite) SetConfiguration(ctx context.Context, data entity.Configuration) (string, error) {
	dataMap, err := data.MapStringInterface()
	if err != nil {
		internalerrors.StackTrace(err)
		return "", err
	}
	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	lastId, err := r.db.DB.InsertOne(r.dbName, doc)
	if err != nil {
		internalerrors.StackTrace(err)
		return "", err
	}

	return lastId, nil
}

func (r *RepositoryApplicationConfigurationWrite) DeleteConfiguration(ctx context.Context, filter entity.FilterConfiguration) error {
	err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()
	if err != nil {
		internalerrors.StackTrace(err)
	}

	return err
}

func (r *RepositoryApplicationConfigurationWrite) UpdateConfiguration(ctx context.Context, data entity.Configuration) error {
	dataMap, err := data.MapStringInterface()
	if err != nil {
		internalerrors.StackTrace(err)
		return err
	}

	err = r.db.DB.
		Query(r.dbName).
		Where(data.FilterConfiguration().Filter()).
		Update(dataMap)
	if err != nil {
		internalerrors.StackTrace(err)
		return err
	}

	return nil
}

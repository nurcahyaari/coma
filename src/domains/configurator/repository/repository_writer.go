package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/configurator/model"
	"github.com/ostafen/clover"
	"github.com/rs/zerolog/log"
)

type RepositoryWriter interface {
	SetConfiguration(ctx context.Context, data model.Configuration) (string, error)
	DeleteConfiguration(ctx context.Context, filter model.FilterConfiguration) error
	UpdateConfiguration(ctx context.Context, data model.Configuration) error
}

type RepositoryWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryWriter(db *database.Clover, name string) RepositoryWriter {
	return &RepositoryWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryWrite) SetConfiguration(ctx context.Context, data model.Configuration) (string, error) {
	dataMap, err := data.MapStringInterface()
	if err != nil {
		return "", err
	}
	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	lastId, err := r.db.DB.InsertOne(r.dbName, doc)
	if err != nil {
		return "", err
	}

	return lastId, nil
}

func (r *RepositoryWrite) DeleteConfiguration(ctx context.Context, filter model.FilterConfiguration) error {
	err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()

	return err
}

func (r *RepositoryWrite) UpdateConfiguration(ctx context.Context, data model.Configuration) error {
	dataMap, err := data.MapStringInterface()
	if err != nil {
		log.Error().Err(err).
			Msg("[UpdateConfiguration] error on create map string interface")
		return err
	}

	err = r.db.DB.
		Query(r.dbName).
		Where(data.FilterConfiguration().Filter()).
		Update(dataMap)
	if err != nil {
		log.Error().Err(err).
			Msg("[UpdateConfiguration] error on delete configuration")
		return err
	}

	return nil
}

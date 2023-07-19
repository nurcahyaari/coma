package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/entity"
	"github.com/coma/coma/src/domains/repository"
	"github.com/ostafen/clover"
	"github.com/rs/zerolog/log"
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

func (r *RepositoryApplicationConfigurationWrite) DeleteConfiguration(ctx context.Context, filter entity.FilterConfiguration) error {
	err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()

	return err
}

func (r *RepositoryApplicationConfigurationWrite) UpdateConfiguration(ctx context.Context, data entity.Configuration) error {
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

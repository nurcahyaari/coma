package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
)

type RepositoryApplicationConfigurationReader interface {
	FindClientConfiguration(ctx context.Context, filter model.FilterConfiguration) (model.Configurations, error)
}

type RepositoryApplicationConfigurationRead struct {
	dbName string
	db     *database.Clover
}

func NewApplicationConfigurationRepositoryReader(db *database.Clover, name string) RepositoryApplicationConfigurationReader {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationConfigurationRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationConfigurationRead) FindClientConfiguration(ctx context.Context, filter model.FilterConfiguration) (model.Configurations, error) {
	var configurations model.Configurations

	docs, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		configuration := model.Configuration{}
		err := doc.Unmarshal(&configuration)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, configuration)
	}

	return configurations, nil
}

package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/entity"
	"github.com/coma/coma/src/domains/repository"
)

type RepositoryApplicationConfigurationRead struct {
	dbName string
	db     *database.Clover
}

func NewApplicationConfigurationRepositoryReader(db *database.Clover, name string) repository.RepositoryApplicationConfigurationReader {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationConfigurationRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationConfigurationRead) FindClientConfiguration(ctx context.Context, filter entity.FilterConfiguration) (entity.Configurations, error) {
	var configurations entity.Configurations

	docs, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		configuration := entity.Configuration{}
		err := doc.Unmarshal(&configuration)
		if err != nil {
			return nil, err
		}
		configurations = append(configurations, configuration)
	}

	return configurations, nil
}

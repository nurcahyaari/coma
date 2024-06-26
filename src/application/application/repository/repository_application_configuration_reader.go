package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
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
		internalerrors.StackTrace(err)
		return nil, err
	}

	for _, doc := range docs {
		configuration := entity.Configuration{}
		err := doc.Unmarshal(&configuration)
		if err != nil {
			internalerrors.StackTrace(err)
			return nil, err
		}
		configurations = append(configurations, configuration)
	}

	return configurations, nil
}

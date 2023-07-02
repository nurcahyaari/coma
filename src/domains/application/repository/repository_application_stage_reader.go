package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
)

type RepositoryApplicationStageReader interface {
	FindStages(ctx context.Context, filter model.FilterApplicationStage) (model.ApplicationStages, error)
}

type RepositoryApplicationStageRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationStageReader(db *database.Clover, name string) RepositoryApplicationStageReader {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationStageRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationStageRead) FindStages(ctx context.Context, filter model.FilterApplicationStage) (model.ApplicationStages, error) {
	var applicationStages model.ApplicationStages

	docs, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		applicationStage := model.ApplicationStage{}
		err := doc.Unmarshal(&applicationStage)
		if err != nil {
			return nil, err
		}
		applicationStages = append(applicationStages, applicationStage)
	}

	return applicationStages, nil
}

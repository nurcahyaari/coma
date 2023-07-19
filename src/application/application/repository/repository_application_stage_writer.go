package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryApplicationStageWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationStageWriter(db *database.Clover, name string) repository.RepositoryApplicationStageWriter {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationStageWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationStageWrite) CreateOrSaveStage(ctx context.Context, data entity.ApplicationStage) error {
	dataMap, err := data.MapStringInterface()
	if err != nil {
		return err
	}

	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	_, err = r.db.DB.InsertOne(r.dbName, doc)
	if err != nil {
		return err
	}

	return nil
}

func (r *RepositoryApplicationStageWrite) DeleteStage(ctx context.Context, filter entity.FilterApplicationStage) error {
	return r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()
}

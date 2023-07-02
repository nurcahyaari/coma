package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
	"github.com/ostafen/clover"
)

type RepositoryApplicationStageWriter interface {
	CreateStage(ctx context.Context, data model.ApplicationStage) error
	DeleteStage(ctx context.Context, filter model.FilterApplicationStage) error
}

type RepositoryApplicationStageWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationStageWriter(db *database.Clover, name string) RepositoryApplicationStageWriter {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationStageWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationStageWrite) CreateStage(ctx context.Context, data model.ApplicationStage) error {
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

func (r *RepositoryApplicationStageWrite) DeleteStage(ctx context.Context, filter model.FilterApplicationStage) error {
	return r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()
}

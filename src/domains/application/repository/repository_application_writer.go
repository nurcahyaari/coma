package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
	"github.com/ostafen/clover"
)

type RepositoryApplicationWriter interface {
	CreateApplication(ctx context.Context, data model.Application) error
	DeleteApplication(ctx context.Context, filter model.FilterApplication) error
}

type RepositoryApplicationWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationWriter(db *database.Clover, name string) RepositoryApplicationWriter {
	return &RepositoryApplicationWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationWrite) CreateApplication(ctx context.Context, data model.Application) error {
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

func (r *RepositoryApplicationWrite) DeleteApplication(ctx context.Context, filter model.FilterApplication) error {
	return r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()
}

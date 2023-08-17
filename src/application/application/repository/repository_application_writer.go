package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryApplicationWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationWriter(db *database.Clover, name string) repository.RepositoryApplicationWriter {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationWrite) CreateApplication(ctx context.Context, data entity.Application) error {
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

func (r *RepositoryApplicationWrite) DeleteApplication(ctx context.Context, filter entity.FilterApplication) error {
	return r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		Delete()
}

package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
	"github.com/ostafen/clover"
)

type RepositoryApplicationKeyWriter interface {
	CreateOrSaveApplicationKey(ctx context.Context, data model.ApplicationKey) error
}

type RepositoryApplicationKeyWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationKeyWriter(db *database.Clover, name string) RepositoryApplicationKeyWriter {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationKeyWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationKeyWrite) CreateOrSaveApplicationKey(ctx context.Context, data model.ApplicationKey) error {
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

package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryApplicationKeyWrite struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationKeyWriter(db *database.Clover, name string) repository.RepositoryApplicationKeyWriter {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationKeyWrite{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationKeyWrite) CreateOrSaveApplicationKey(ctx context.Context, data entity.ApplicationKey) error {
	dataMap, err := data.MapStringInterface()
	if err != nil {
		internalerrors.StackTrace(err)
		return err
	}

	doc := clover.NewDocument()
	doc.SetAll(dataMap)

	_, err = r.db.DB.InsertOne(r.dbName, doc)
	if err != nil {
		internalerrors.StackTrace(err)
		return err
	}

	return nil
}

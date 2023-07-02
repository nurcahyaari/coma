package repository

import (
	"context"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
)

type RepositoryApplicationReader interface {
	FindApplications(ctx context.Context, filter model.FilterApplication) (model.Applications, error)
}

type RepositoryApplicationRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationReader(db *database.Clover, name string) RepositoryApplicationReader {
	return &RepositoryApplicationRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationRead) FindApplications(ctx context.Context, filter model.FilterApplication) (model.Applications, error) {
	var applications model.Applications

	docs, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		application := model.Application{}
		err := doc.Unmarshal(&application)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

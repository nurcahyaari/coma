package repository

import (
	"context"
	"errors"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/application/model"
)

type RepositoryApplicationReader interface {
	FindApplication(ctx context.Context, filter model.FilterApplication) (model.Application, error)
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

func (r *RepositoryApplicationRead) FindApplication(ctx context.Context, filter model.FilterApplication) (model.Application, error) {
	applications, err := r.FindApplications(ctx, filter)
	if err != nil {
		return model.Application{}, err
	}
	if len(applications) == 0 {
		return model.Application{}, errors.New("err: data is not found")
	}

	return applications[0], nil
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

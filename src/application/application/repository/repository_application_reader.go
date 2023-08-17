package repository

import (
	"context"
	"errors"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domain/entity"
	"github.com/coma/coma/src/domain/repository"
)

type RepositoryApplicationRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryApplicationReader(db *database.Clover, name string) repository.RepositoryApplicationReader {
	db.DB.CreateCollection(name)
	return &RepositoryApplicationRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryApplicationRead) FindApplication(ctx context.Context, filter entity.FilterApplication) (entity.Application, error) {
	applications, err := r.FindApplications(ctx, filter)
	if err != nil {
		return entity.Application{}, err
	}
	if len(applications) == 0 {
		return entity.Application{}, errors.New("err: application is not found")
	}

	return applications[0], nil
}

func (r *RepositoryApplicationRead) FindApplications(ctx context.Context, filter entity.FilterApplication) (entity.Applications, error) {
	var applications entity.Applications

	docs, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		return nil, err
	}

	for _, doc := range docs {
		application := entity.Application{}
		err := doc.Unmarshal(&application)
		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

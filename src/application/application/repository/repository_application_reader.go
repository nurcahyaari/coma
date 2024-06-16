package repository

import (
	"context"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
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

func (r *RepositoryApplicationRead) FindApplication(ctx context.Context, filter entity.FilterApplication) (entity.Application, bool, error) {
	applications, err := r.FindApplications(ctx, filter)
	if err != nil {
		internalerrors.StackTrace(err)
		return entity.Application{}, false, err
	}
	if len(applications) == 0 {
		return entity.Application{}, false, nil
	}

	return applications[0], true, nil
}

func (r *RepositoryApplicationRead) FindApplications(ctx context.Context, filter entity.FilterApplication) (entity.Applications, error) {
	var applications entity.Applications

	docs, err := r.db.DB.
		Query(r.dbName).
		Where(filter.Filter()).
		FindAll()
	if err != nil {
		internalerrors.StackTrace(err)
		return nil, err
	}

	for _, doc := range docs {
		application := entity.Application{}
		err := doc.Unmarshal(&application)
		if err != nil {
			internalerrors.StackTrace(err)
			return nil, err
		}
		applications = append(applications, application)
	}

	return applications, nil
}

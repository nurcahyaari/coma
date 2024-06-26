package repository

import (
	"context"
	"errors"

	"github.com/nurcahyaari/coma/infrastructure/database"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/nurcahyaari/coma/src/domain/repository"
	"github.com/ostafen/clover"
)

type RepositoryUserRead struct {
	dbName string
	db     *database.Clover
}

func NewRepositoryUserReader(db *database.Clover, name string) repository.RepositoryUserAuthReader {
	db.DB.CreateCollection(name)
	return &RepositoryUserRead{
		db:     db,
		dbName: name,
	}
}

func (r *RepositoryUserRead) FindTokenBy(ctx context.Context, filter entity.FilterUserAuth) (*entity.UserAuth, error) {
	var (
		resp *entity.UserAuth
	)

	criteria := filter.Filter()

	if criteria == nil {
		err := errors.New("err: filter cannot be nulled")
		internalerrors.StackTrace(err)
		return nil, err
	}

	doc, err := r.db.DB.Query(r.dbName).
		Where(criteria).
		Sort(clover.SortOption{
			Field:     "accessTokenExpiredAt",
			Direction: -1,
		}).
		FindFirst()
	if err != nil {
		internalerrors.StackTrace(err)
		return nil, err
	}
	if doc == nil {
		return resp, nil
	}

	if err := doc.Unmarshal(&resp); err != nil {
		internalerrors.StackTrace(err)
		return nil, err
	}

	return resp, nil
}

package repository

import (
	"context"
	"fmt"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/auth/model"
	"github.com/ostafen/clover"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryReader
type RepositoryReader interface {
	FindTokenById(ctx context.Context, id int64) (model.Apikey, error)
	FindTokenByToken(ctx context.Context, token string) (model.Apikey, error)
}

type RepositoryRead struct {
	db *database.Clover
}

//counterfeiter:generate . RepositoryRead
func NewRepositoryReader(db *database.Clover) RepositoryReader {
	db.DB.CreateCollection("apikey")
	return &RepositoryRead{
		db: db,
	}
}

func (r *RepositoryRead) FindTokenById(ctx context.Context, id int64) (model.Apikey, error) {
	var apiKey model.Apikey
	doc, err := r.db.DB.Query("apikey").FindById(fmt.Sprintf("%d", id))
	if err != nil {
		return apiKey, err
	}

	err = doc.Unmarshal(&apiKey)
	if err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func (r *RepositoryRead) FindTokenByToken(ctx context.Context, token string) (model.Apikey, error) {
	var apiKey model.Apikey

	criteria := clover.Field("key").
		Eq(token)

	doc, err := r.db.DB.Query("apikey").
		Where(criteria).
		FindFirst()
	if err != nil {
		return apiKey, err
	}

	err = doc.Unmarshal(&apiKey)
	if err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

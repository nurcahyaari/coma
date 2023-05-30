package repository

import (
	"context"
	"fmt"

	"github.com/coma/coma/infrastructure/database"
	"github.com/coma/coma/src/domains/auth/model"
	"github.com/ostafen/clover"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryReader
type ApiKeyRepositoryReader interface {
	FindTokenById(ctx context.Context, id int64) (model.Apikey, error)
	FindTokenByToken(ctx context.Context, token string) (model.Apikey, error)
}

type ApiKeyRepository struct {
	DB *database.Clover
}

//counterfeiter:generate . Repository
func NewApiKey(db *database.Clover) ApiKeyRepositoryReader {
	db.DB.CreateCollection("apikey")
	return &ApiKeyRepository{
		DB: db,
	}
}

func (r *ApiKeyRepository) FindTokenById(ctx context.Context, id int64) (model.Apikey, error) {
	var apiKey model.Apikey
	doc, err := r.DB.DB.Query("apikey").FindById(fmt.Sprintf("%d", id))
	if err != nil {
		return apiKey, err
	}

	err = doc.Unmarshal(&apiKey)
	if err != nil {
		return apiKey, err
	}

	return apiKey, nil
}

func (r *ApiKeyRepository) FindTokenByToken(ctx context.Context, token string) (model.Apikey, error) {
	var apiKey model.Apikey

	criteria := clover.Field("key").
		Eq(token)

	doc, err := r.DB.DB.Query("apikey").
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

package repository

import (
	"context"

	"github.com/coma/coma/src/domain/entity"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryAuthReader
type RepositoryAuthReader interface {
	FindTokenById(ctx context.Context, id int64) (entity.Apikey, error)
	FindTokenByToken(ctx context.Context, token string) (entity.Apikey, error)
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryAuthWriter
type RepositoryAuthWriter interface{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . AuthRepositorier
type AuthRepositorier interface {
	NewRepositoryReader() RepositoryAuthReader
	NewRepositoryWriter() RepositoryAuthWriter
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryOAuthReader
type RepositoryOAuthReader interface{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . RepositoryOAuthWriter
type RepositoryOAuthWriter interface{}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 . OauthRepositorier
type OauthRepositorier interface {
	NewRepositoryOAuthReader() RepositoryOAuthReader
	NewRepositoryOAuthWriter() RepositoryOAuthWriter
}

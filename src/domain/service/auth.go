package service

import (
	"context"

	"github.com/nurcahyaari/coma/src/application/auth/dto"
)

type AuthServicer interface {
	GenerateToken(context.Context, dto.RequestGenerateToken) (dto.ResponseGenerateToken, error)
	ValidateToken(context.Context, dto.RequestValidateToken) (dto.ResponseValidateKey, error)
	ExtractToken(context.Context, dto.RequestValidateToken) (dto.ResponseExtractedToken, error)
}

type LocalUserAuthServicer interface {
	AuthServicer
	ValidateUserScope(context.Context, dto.RequestUserScopeValidation) (dto.ResponseValidateKey, error)
	ValidateUserApplicationScope(context.Context, dto.RequestUserApplicationScopeValidation) (dto.ResponseValidateKey, error)
}

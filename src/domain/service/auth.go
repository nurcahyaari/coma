package service

import (
	"context"

	"github.com/coma/coma/src/application/auth/dto"
)

type AuthServicer interface {
	ValidateToken(context.Context, dto.RequestValidateToken) (dto.ResponseValidateKey, error)
}

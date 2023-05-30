package service

import (
	"context"

	"github.com/coma/coma/src/domains/auth/dto"
)

type AuthServicer interface {
	ValidateToken(context.Context, dto.RequestAuthValidate) (dto.ResponseValidateKey, error)
}

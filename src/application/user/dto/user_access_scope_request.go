package dto

import (
	"net/http"

	internalerror "github.com/coma/coma/internal/x/errors"
	"github.com/coma/coma/src/domain/entity"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RequestCreateUserApplicationScope struct {
	UserId        string                    `json:"userId"`
	ApplicationId string                    `json:"applicationId"`
	Rbac          *UserApplicationScopeRbac `json:"rbac"`
}

func (r RequestCreateUserApplicationScope) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.UserId, validation.Required),
		validation.Field(&r.ApplicationId, validation.Required),
	)

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest))
}

func (r RequestCreateUserApplicationScope) UserApplicationScope() entity.UserApplicationScope {
	id := uuid.New()
	userApplicationScope := entity.UserApplicationScope{
		Id:            id.String(),
		ApplicationId: r.ApplicationId,
		UserId:        r.UserId,
		Rbac:          &entity.UserApplicationScopeRbac{},
	}

	if r.Rbac != nil {
		userApplicationScope.Rbac = r.Rbac.UserApplicationScopeRbac()
	}

	return userApplicationScope
}

type RequestFindUserApplicationScope struct {
	UserId        string `json:"userId"`
	ApplicationId string `json:"applicationId"`
}

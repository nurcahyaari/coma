package dto

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
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

	return internalerror.New(err,
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

package dto

import (
	"github.com/coma/coma/src/domain/entity"
	"github.com/google/uuid"
)

type RequestCreateUserApplicationScope struct {
	UserId        string                    `json:"userId"`
	ApplicationId string                    `json:"applicationId"`
	StageId       string                    `json:"stageId"`
	Rbac          *UserApplicationScopeRbac `json:"rbac"`
}

func (r RequestCreateUserApplicationScope) UserApplicationScope() entity.UserApplicationScope {
	id := uuid.New()
	userApplicationScope := entity.UserApplicationScope{
		Id:            id.String(),
		ApplicationId: r.ApplicationId,
		StageId:       r.StageId,
		Rbac:          &entity.UserApplicationScopeRbac{},
	}

	if r.Rbac != nil {
		userApplicationScope.Rbac = r.Rbac.UserApplicationScopeRbac()
	}

	return userApplicationScope
}

type RequestFindUserApplicationScope struct {
	UserId string `json:"userId"`
}

package dto

import (
	"github.com/coma/coma/src/domain/entity"
	"github.com/google/uuid"
)

type RequestCreateUserAccessScope struct {
	UserId        string               `json:"userId"`
	ApplicationId string               `json:"applicationId"`
	StageId       string               `json:"stageId"`
	Rbac          *UserAccessScopeRbac `json:"rbac"`
}

func (r RequestCreateUserAccessScope) UserAccessScope() entity.UserAccessScope {
	id := uuid.New()
	userAccessScope := entity.UserAccessScope{
		Id:            id.String(),
		ApplicationId: r.ApplicationId,
		StageId:       r.StageId,
	}

	if r.Rbac != nil {
		userAccessScope.Rbac = r.Rbac.UserAccessScopeRbac()
	}

	return userAccessScope
}

type RequestFindUserAccessScope struct {
	UserId string `json:"userId"`
}

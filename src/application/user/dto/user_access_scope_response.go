package dto

import "github.com/coma/coma/src/domain/entity"

type ResponseUserAccessScope struct {
	UserId string               `json:"userId"`
	Rbac   *UserAccessScopeRbac `json:"rbac"`
}

type UserAccessScopeRbac struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

func (r UserAccessScopeRbac) UserAccessScopeRbac() *entity.UserAccessScopeRbac {
	return &entity.UserAccessScopeRbac{
		Create: r.Create,
		Read:   r.Read,
		Update: r.Update,
		Delete: r.Delete,
	}
}

func NewResponseUserAccessScope(userAccessScope entity.UserAccessScope) ResponseUserAccessScope {
	return ResponseUserAccessScope{
		UserId: userAccessScope.UserId,
		Rbac: &UserAccessScopeRbac{
			Create: userAccessScope.Rbac.Create,
			Read:   userAccessScope.Rbac.Read,
			Update: userAccessScope.Rbac.Update,
			Delete: userAccessScope.Rbac.Delete,
		},
	}
}

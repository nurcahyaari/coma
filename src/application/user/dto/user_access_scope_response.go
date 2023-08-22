package dto

import "github.com/coma/coma/src/domain/entity"

type ResponseUserApplicationScope struct {
	UserId string                    `json:"userId"`
	Rbac   *UserApplicationScopeRbac `json:"rbac"`
}

type UserApplicationScopeRbac struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

func (r UserApplicationScopeRbac) UserApplicationScopeRbac() *entity.UserApplicationScopeRbac {
	return &entity.UserApplicationScopeRbac{
		Create: r.Create,
		Read:   r.Read,
		Update: r.Update,
		Delete: r.Delete,
	}
}

func NewResponseUserApplicationScope(userApplicationScope entity.UserApplicationScope) ResponseUserApplicationScope {
	return ResponseUserApplicationScope{
		UserId: userApplicationScope.UserId,
		Rbac: &UserApplicationScopeRbac{
			Create: userApplicationScope.Rbac.Create,
			Read:   userApplicationScope.Rbac.Read,
			Update: userApplicationScope.Rbac.Update,
			Delete: userApplicationScope.Rbac.Delete,
		},
	}
}

type ResponseUserApplicationsScope []ResponseUserApplicationScope

func NewResponseUserApplicationsScope(userApplicationsScope entity.UserApplicationsScope) ResponseUserApplicationsScope {
	resp := make(ResponseUserApplicationsScope, 0)

	for _, r := range userApplicationsScope {
		resp = append(resp, NewResponseUserApplicationScope(r))
	}

	return resp
}

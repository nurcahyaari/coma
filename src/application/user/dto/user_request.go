package dto

import (
	"github.com/coma/coma/src/domain/entity"
	"github.com/google/uuid"
)

type RequestUser struct {
	Id       string          `json:"-"`
	Username string          `json:"username"`
	UserType entity.UserType `json:"-"`
}

func (r RequestUser) User() entity.User {
	return entity.User{
		Id:       r.Id,
		Username: r.Username,
	}
}

type RequestUsers struct {
	Page int
	Size int
}

type UserRbac struct {
	Create bool `json:"create"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type RequestCreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RequestCreateUserNonRoot struct {
	RequestCreateUser
	Rbac UserRbac `json:"rbac"`
}

func (r RequestCreateUserNonRoot) User() entity.User {
	uid := uuid.New()
	return entity.User{
		Id:       uid.String(),
		Username: r.Username,
		Password: r.Password,
		UserType: entity.UserTypeUser,
		Rbac: &entity.UserRbac{
			Create: r.Rbac.Create,
			Delete: r.Rbac.Update,
			Update: r.Rbac.Delete,
		},
	}
}

func (r RequestCreateUser) UserRoot() entity.User {
	uid := uuid.New()
	return entity.User{
		Id:       uid.String(),
		Username: r.Username,
		Password: r.Password,
		UserType: entity.UserTypeRoot,
	}
}

type RequestPatchUserPassword struct {
	Id       string `json:"-"`
	Passowrd string `json:"password"`
}

package dto

import (
	"github.com/coma/coma/src/domain/entity"
	"github.com/google/uuid"
)

type RequestUser struct {
	Id       string `json:"-"`
	Username string `json:"username"`
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

type RequestCreateUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r RequestCreateUser) User() entity.User {
	uid := uuid.New()
	return entity.User{
		Id:       uid.String(),
		Username: r.Username,
		Password: r.Password,
	}
}

type RequestPatchUserPassword struct {
	Id       string `json:"-"`
	Passowrd string `json:"password"`
}

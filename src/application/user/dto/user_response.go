package dto

import "github.com/coma/coma/src/domain/entity"

type ResponseUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func NewResponseUser(user entity.User) ResponseUser {
	return ResponseUser{
		Id:       user.Id,
		Username: user.Username,
	}
}

type ResponseUsers []ResponseUser

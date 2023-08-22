package dto

import "github.com/coma/coma/src/domain/entity"

type ResponseUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	UserType string `json:"userType"`
}

func NewResponseUser(user entity.User) ResponseUser {
	return ResponseUser{
		Id:       user.Id,
		Username: user.Username,
		UserType: string(user.UserType),
	}
}

type ResponseUsers []ResponseUser

func NewResponseUsers(users entity.Users) ResponseUsers {
	resp := make(ResponseUsers, 0)

	for _, u := range users {
		resp = append(resp, NewResponseUser(u))
	}

	return resp
}

package dto

import "github.com/coma/coma/src/domain/entity"

type RequestValidateToken struct {
	Token     string
	TokenType entity.TokenType
}

type RequestGenerateToken struct {
	// key can be username
	Key string `json:"key"`
	// secret can be user password
	Secret string `json:"secret"`
}

type RequestUserApplicationScopeValidation struct {
	UserId string
	Method string
}

type RequestUserScopeValidation struct {
	UserId string
	Method string
}

type ResponseGenerateToken struct {
	AccessToken     string `json:"accessToken"`
	AccessTokenExp  string `json:"accessTokenExp"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp string `json:"refreshTokenExp"`
}

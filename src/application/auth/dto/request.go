package dto

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
)

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

func (r RequestGenerateToken) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Key, validation.Required),
		validation.Field(&r.Secret, validation.Required),
	)

	if err == nil {
		return nil
	}

	return internalerror.New(err,
		internalerror.SetErrorCode(http.StatusBadRequest))
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

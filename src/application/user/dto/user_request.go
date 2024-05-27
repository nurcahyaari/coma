package dto

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
)

type RequestUser struct {
	Id       string          `json:"-"`
	Username string          `json:"username"`
	UserType entity.UserType `json:"-"`
}

func (r RequestUser) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Id, validation.Required),
		validation.Field(&r.Username, validation.Required),
	)

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest))
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

func (r RequestCreateUser) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest))
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

func (r RequestPatchUserPassword) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Id, validation.Required),
		validation.Field(&r.Passowrd, validation.Required),
	)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest))
}

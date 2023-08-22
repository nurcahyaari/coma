package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	userdto "github.com/coma/coma/src/application/user/dto"
	"github.com/go-chi/chi/v5"
)

// FindUser find user
// @Summary find user
// @Description find user
// @Param authorization header string true "User accessToken"
// @Param id path string true "user id"
// @Tags Users
// @Produce json
// @Router /v1/users/{id} [GET]
func (h *HttpHandle) FindUser(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestUser{
		Id: chi.URLParam(r, "id"),
	}

	resp, err := h.userSvc.FindUser(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()),
			response.SetHttpCode[string](errCustom.ErrCode))
		return
	}

	response.Json[userdto.ResponseUser](w,
		response.SetMessage[userdto.ResponseUser]("success"),
		response.SetData[userdto.ResponseUser](resp))
}

// FindUsers find users
// @Summary find users
// @Description find users
// @Param authorization header string true "User accessToken"
// @Tags Users
// @Produce json
// @Router /v1/users [GET]
func (h *HttpHandle) FindUsers(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestUsers{}

	resp, err := h.userSvc.FindUsers(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()),
			response.SetHttpCode[string](errCustom.ErrCode))
		return
	}

	response.Json[userdto.ResponseUsers](w,
		response.SetMessage[userdto.ResponseUsers]("success"),
		response.SetData[userdto.ResponseUsers](resp))
}

// CreateUser set new users
// @Summary set new users
// @Description set new users
// @Param authorization header string true "User accessToken"
// @Param RequestCreateUserNonRoot body userdto.RequestCreateUserNonRoot true "create new user"
// @Tags Users
// @Produce json
// @Router /v1/users [POST]
func (h *HttpHandle) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestCreateUserNonRoot{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()),
			response.SetHttpCode[string](errCustom.ErrCode))
		return
	}

	// TODO: make validation
	resp, err := h.userSvc.CreateUser(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()),
			response.SetHttpCode[string](errCustom.ErrCode))
		return
	}

	response.Json[userdto.ResponseUser](w,
		response.SetMessage[userdto.ResponseUser]("success"),
		response.SetData[userdto.ResponseUser](resp))
}

// CreateUserRoot set new users as root access
// @Summary set new users as root access
// @Description set new users as root access
// @Param RequestCreateUser body userdto.RequestCreateUser true "create new user"
// @Tags Users
// @Produce json
// @Router /v1/users/root [POST]
func (h *HttpHandle) CreateUserRoot(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestCreateUser{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	// TODO: make validation
	resp, err := h.userSvc.CreateRootUser(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()),
			response.SetHttpCode[string](errCustom.ErrCode))
		return
	}

	response.Json[userdto.ResponseUser](w,
		response.SetMessage[userdto.ResponseUser]("success"),
		response.SetData[userdto.ResponseUser](resp))
}

// DeleteUser delete users
// @Summary delete users
// @Description delete users
// @Param authorization header string true "User accessToken"
// @Param id path string true "user id"
// @Tags Users
// @Produce json
// @Router /v1/users/{id} [DELETE]
func (h *HttpHandle) DeleteUser(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestUser{
		Id: chi.URLParam(r, "id"),
	}

	// TODO: make validation
	err := h.userSvc.DeleteUser(r.Context(), req)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

// UpdateUser update users
// @Summary update users
// @Description update users
// @Param authorization header string true "User accessToken"
// @Param id path string true "user id"
// @Param RequestUser body userdto.RequestUser true "update user"
// @Tags Users
// @Produce json
// @Router /v1/users/{id} [PUT]
func (h *HttpHandle) UpdateUser(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestUser{
		Id: chi.URLParam(r, "id"),
	}

	// TODO: make validation
	resp, err := h.userSvc.UpdateUser(r.Context(), req)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[userdto.ResponseUser](w,
		response.SetMessage[userdto.ResponseUser]("success"),
		response.SetData[userdto.ResponseUser](resp))
}

// UpdateUserPassword update users password
// @Summary update users password
// @Description update users password
// @Param authorization header string true "User accessToken"
// @Param id path string true "user id"
// @Param RequestPatchUserPassword body userdto.RequestPatchUserPassword true "update user password"
// @Tags Users
// @Produce json
// @Router /v1/users/password/{id} [PATCH]
func (h *HttpHandle) UpdateUserPassword(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestPatchUserPassword{
		Id: chi.URLParam(r, "id"),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	// TODO: make validation
	err := h.userSvc.PatchUserPassword(r.Context(), req)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

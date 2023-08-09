package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	userdto "github.com/coma/coma/src/application/user/dto"
	"github.com/go-chi/chi/v5"
)

// FindUser find user
// @Summary find user
// @Description find user
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
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[userdto.ResponseUser](w,
		response.SetMessage[userdto.ResponseUser]("success"),
		response.SetData[userdto.ResponseUser](resp))
}

// CreateUser set new users
// @Summary set new users
// @Description set new users
// @Param RequestCreateUser body userdto.RequestCreateUser true "create new user"
// @Tags Users
// @Produce json
// @Router /v1/users [POST]
func (h *HttpHandle) CreateUser(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestCreateUser{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	// TODO: make validation
	resp, err := h.userSvc.CreateUser(r.Context(), req)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[userdto.ResponseUser](w,
		response.SetMessage[userdto.ResponseUser]("success"),
		response.SetData[userdto.ResponseUser](resp))
}

// DeleteUser delete users
// @Summary delete users
// @Description delete users
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
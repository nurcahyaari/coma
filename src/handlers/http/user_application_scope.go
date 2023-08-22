package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	userdto "github.com/coma/coma/src/application/user/dto"
)

// FindUserApplicationScope find user scope to application
// @Summary find user scope to application
// @Description find user scope to application
// @Param authorization header string true "User accessToken"
// @Tags UserApplicationScope
// @Produce json
// @Router /v1/users/application/scope [GET]
func (h *HttpHandle) FindUserApplicationsScope(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestFindUserApplicationScope{}

	resp, err := h.userApplicationScopeSvc.FindUserApplicationsScope(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()),
			response.SetHttpCode[string](errCustom.ErrCode))
		return
	}

	response.Json[userdto.ResponseUserApplicationsScope](w,
		response.SetMessage[userdto.ResponseUserApplicationsScope]("success"),
		response.SetData[userdto.ResponseUserApplicationsScope](resp))
}

// CreateOrUpdateUserApplicationScope create or update user application scope
// @Summary create or update user application scope
// @Description create or update user application scope
// @Param authorization header string true "User accessToken"
// @Param RequestCreateUserApplicationScope body userdto.RequestCreateUserApplicationScope true "create or update application"
// @Tags UserApplicationScope
// @Produce json
// @Router /v1/users/application/scope [POST]
func (h *HttpHandle) CreateOrUpdateUserApplicationScope(w http.ResponseWriter, r *http.Request) {
	req := userdto.RequestCreateUserApplicationScope{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	err := h.userApplicationScopeSvc.UpsetUserApplicationScope(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[any](w,
			response.SetErr[any](errCustom.ErrorAsObject()),
			response.SetHttpCode[any](errCustom.ErrCode))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

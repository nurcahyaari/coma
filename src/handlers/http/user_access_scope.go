package http

import (
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	userdto "github.com/coma/coma/src/application/user/dto"
	"github.com/go-chi/chi/v5"
)

// FindUserApplicationScope find user access scope to application
// @Summary find user access scope to application
// @Description find user access scope to application
// @Param authorization header string true "User accessToken"
// @Tags Users
// @Produce json
// @Router /v1/users/application/access [GET]
func (h *HttpHandle) FindUserApplicationScope(w http.ResponseWriter, r *http.Request) {
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

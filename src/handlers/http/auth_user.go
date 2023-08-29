package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/application/auth/dto"
)

// AuthUserLogin login auth user by username and password
// @Summary login auth user by username and password
// @Description login auth user by username and password
// @Param RequestGenerateToken body dto.RequestGenerateToken true "login"
// @Tags Auth
// @Produce json
// @Router /v1/auth/user/login [POST]
func (h *HttpHandle) AuthUserLogin(w http.ResponseWriter, r *http.Request) {
	req := dto.RequestGenerateToken{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	if err := req.Validate(); err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()))
		return
	}

	resp, err := h.authSvc.GenerateToken(r.Context(), req)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()))
		return
	}

	response.Json[dto.ResponseGenerateToken](w,
		response.SetMessage[dto.ResponseGenerateToken]("success"),
		response.SetData[dto.ResponseGenerateToken](resp))
}

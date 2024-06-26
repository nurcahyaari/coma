package http

import (
	"encoding/json"
	"net/http"

	"github.com/nurcahyaari/coma/internal/protocols/http/response"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	applicationdto "github.com/nurcahyaari/coma/src/application/application/dto"
)

// FindApplicationKey get key
// @Summary get key
// @Security comaStandardAuth
// @Description get key
// @Param applicationId query string false "<Application Id>"
// @Tags Key
// @Produce json
// @Router /v1/keys [GET]
func (h *HttpHandle) FindApplicationKey(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestFindApplicationKey{
		ApplicationId: r.FormValue("applicationId"),
	}

	resp, err := h.applicationKeySvc.FindApplicationKey(r.Context(), request)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[any](w,
			response.SetErr[any](errCustom.ErrorAsObject()),
			response.SetHttpCode[any](errCustom.ErrCode))
		return
	}

	response.Json[applicationdto.ResponseFindApplicationKey](w,
		response.SetMessage[applicationdto.ResponseFindApplicationKey]("success"),
		response.SetData[applicationdto.ResponseFindApplicationKey](resp))
}

// CreateOrUpdateApplicationKey create or update existing key
// @Summary create or update existing key
// @Security comaStandardAuth
// @Description create or update existing key
// @Param RequestCreateApplicationKey body applicationdto.RequestCreateApplicationKey true "create new stages"
// @Tags Key
// @Produce json
// @Router /v1/keys [POST]
func (h *HttpHandle) CreateOrUpdateApplicationKey(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestCreateApplicationKey{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[any](w,
			response.SetErr[any](errCustom.ErrorAsObject()))
		return
	}

	resp, err := h.applicationKeySvc.GenerateOrUpdateApplicationKey(r.Context(), request)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[any](w,
			response.SetErr[any](errCustom.ErrorAsObject()),
			response.SetHttpCode[any](errCustom.ErrCode))
		return
	}

	response.Json[applicationdto.ResponseCreateApplicationKey](w,
		response.SetMessage[applicationdto.ResponseCreateApplicationKey]("success"),
		response.SetData[applicationdto.ResponseCreateApplicationKey](resp))
}

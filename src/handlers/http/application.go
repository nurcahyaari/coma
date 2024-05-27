package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nurcahyaari/coma/internal/protocols/http/response"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	applicationdto "github.com/nurcahyaari/coma/src/application/application/dto"
)

// FindApplications get applications
// @Summary get applications
// @Security comaStandardAuth
// @Description get applications
// @Param applicationName query string false "<application name>"
// @Tags Applications
// @Produce json
// @Router /v1/applications [GET]
func (h *HttpHandle) FindApplications(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestFindApplication{
		Name: r.FormValue("applicationName"),
	}

	resp, err := h.applicationSvc.FindApplications(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[applicationdto.ResponseApplications](w,
		response.SetMessage[applicationdto.ResponseApplications]("success"),
		response.SetData[applicationdto.ResponseApplications](resp))
}

// CreateApplication set new applications
// @Summary set new applications
// @Security comaStandardAuth
// @Description Set new applications
// @Param RequestCreateApplication body applicationdto.RequestCreateApplication true "create new application"
// @Tags Applications
// @Produce json
// @Router /v1/applications [POST]
func (h *HttpHandle) CreateApplication(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestCreateApplication{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[any](w,
			response.SetMessage[any](err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[any](w,
			response.SetErr[any](errCustom.ErrorAsObject()))
		return
	}

	resp, err := h.applicationSvc.CreateApplication(r.Context(), request)
	if err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[any](w,
			response.SetErr[any](errCustom.ErrorAsObject()),
			response.SetHttpCode[any](errCustom.ErrCode))
		return
	}

	response.Json[applicationdto.ResponseApplication](w,
		response.SetMessage[applicationdto.ResponseApplication]("success"),
		response.SetData[applicationdto.ResponseApplication](resp))
}

// DeleteApplications delete application
// @Summary delete application
// @Security comaStandardAuth
// @Description delete application
// @Param applicationId path string true "application id"
// @Tags Applications
// @Produce json
// @Router /v1/applications/{applicationId} [DELETE]
func (h *HttpHandle) DeleteApplications(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestFindApplication{
		Id: chi.URLParam(r, "applicationId"),
	}

	err := h.applicationSvc.DeleteApplication(r.Context(), request)
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

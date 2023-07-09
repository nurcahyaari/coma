package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	internalerrors "github.com/coma/coma/internal/utils/errors"
	applicationdto "github.com/coma/coma/src/domains/application/dto"
	"github.com/go-chi/chi/v5"
)

// FindApplications get applications
// @Summary get applications
// @Description get applications
// @Param stageId query string false "<Stage id>"
// @Param applicationName query string false "<application name>"
// @Tags Applications
// @Produce json
// @Router /v1/applications [GET]
func (h *HttpHandle) FindApplications(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestFindApplication{
		Name:    r.FormValue("applicationName"),
		StageId: r.FormValue("stageId"),
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

// CreateApplicationStages set new applications
// @Summary set new applications
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
// @Description delete application
// @Param applicationId path int true "application id"
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

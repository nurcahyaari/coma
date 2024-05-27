package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/nurcahyaari/coma/internal/protocols/http/response"
	internalerrors "github.com/nurcahyaari/coma/internal/x/errors"
	applicationdto "github.com/nurcahyaari/coma/src/application/application/dto"
)

// GetConfiguration get it's config
// @Summary set new config
// @Security comaStandardAuth
// @Description Set new config
// @Param x-clientkey header string true "<Client Key>"
// @Param viewType query string true "<View Type>" Enums(JSON, schema)
// @Tags Config
// @Produce json
// @Router /v1/configuration [GET]
func (h *HttpHandle) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestGetConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	viewType := r.FormValue("viewType")

	switch viewType {
	case applicationdto.ViewTypeJSON:
		resp, err := h.configurationSvc.GetConfigurationViewTypeJSON(r.Context(), request)
		if err != nil {
			response.Err[string](w,
				response.SetErr[string](err.Error()))
			return
		}

		response.Json[applicationdto.ResponseGetConfigurationViewTypeJSON](w,
			response.SetData[applicationdto.ResponseGetConfigurationViewTypeJSON](resp),
			response.SetMessage[applicationdto.ResponseGetConfigurationViewTypeJSON]("success"))
		return
	default:
		resp, err := h.configurationSvc.GetConfigurationViewTypeSchema(r.Context(), request)
		if err != nil {
			response.Err[string](w,
				response.SetErr[string](err.Error()))
			return
		}

		response.Json[applicationdto.ResponseGetConfigurationsViewTypeSchema](w,
			response.SetData[applicationdto.ResponseGetConfigurationsViewTypeSchema](resp),
			response.SetMessage[applicationdto.ResponseGetConfigurationsViewTypeSchema]("success"))
		return
	}
}

// SetConfiguration set new config
// @Summary set new config
// @Security comaStandardAuth
// @Description Set new config
// @Param x-clientkey header string true "<Client Key>"
// @Param RequestSetConfiguration body applicationdto.RequestSetConfiguration true "create new field of config"
// @Tags Config
// @Produce json
// @Router /v1/configuration [POST]
func (h *HttpHandle) SetConfiguration(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestSetConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

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

	res, err := h.configurationSvc.SetConfiguration(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[applicationdto.ResponseSetConfiguration](w,
		response.SetMessage[applicationdto.ResponseSetConfiguration]("success"),
		response.SetData[applicationdto.ResponseSetConfiguration](res))
}

// UpdateConfiguration update new config
// @Summary update new config
// @Security comaStandardAuth
// @Description update new config
// @Tags Config
// @Param x-clientkey header string true "<Client Key>"
// @Param RequestUpdateConfiguration body applicationdto.RequestUpdateConfiguration true "update data of config"
// @Produce json
// @Router /v1/configuration [PUT]
func (h *HttpHandle) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestUpdateConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()))
		return
	}

	err = h.configurationSvc.UpdateConfiguration(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

// UpsertConfiguration update/set configuration
// @Summary update or set configuration
// @Security comaStandardAuth
// @Description update or set configuration
// @Param x-clientkey header string true "<Client Key>"
// @Param RequestSetConfiguration body applicationdto.RequestSetConfiguration true "create new field of config"
// @Tags Config
// @Produce json
// @Router /v1/configuration/upsert [POST]
func (h *HttpHandle) UpsertConfiguration(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestSetConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	if err := request.Validate(); err != nil {
		errCustom := err.(*internalerrors.Error)
		response.Err[string](w,
			response.SetErr[string](errCustom.Error()))
		return
	}

	err = h.configurationSvc.UpsertConfiguration(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

// DeleteConfiguration delete configuration based on it's id
// @Summary delete a config
// @Security comaStandardAuth
// @Description delete a config
// @Param x-clientkey header string true "<Client Key>"
// @Param id path string true "The config field identifier it's a UUID."
// @Tags Config
// @Produce json
// @Router /v1/configuration/{id} [DELETE]
func (h *HttpHandle) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestDeleteConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
		Id:         chi.URLParam(r, "id"),
	}

	err := h.configurationSvc.DeleteConfiguration(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

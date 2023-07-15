package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	configuratordto "github.com/coma/coma/src/domains/application/dto"
	"github.com/go-chi/chi/v5"
)

// GetConfiguration get it's config
// @Summary set new config
// @Description Set new config
// @Param x-clientkey header string true "<Client Key>"
// @Param viewType query string true "<View Type>" Enums(JSON, schema)
// @Tags Config
// @Produce json
// @Router /v1/configurator [GET]
func (h *HttpHandle) GetConfiguration(w http.ResponseWriter, r *http.Request) {
	request := configuratordto.RequestGetConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	viewType := r.FormValue("viewType")

	switch viewType {
	case configuratordto.ViewTypeJSON:
		resp, err := h.configurationSvc.GetConfigurationViewTypeJSON(r.Context(), request)
		if err != nil {
			response.Err[string](w,
				response.SetErr[string](err.Error()))
			return
		}

		response.Json[configuratordto.ResponseGetConfigurationViewTypeJSON](w,
			response.SetData[configuratordto.ResponseGetConfigurationViewTypeJSON](resp),
			response.SetMessage[configuratordto.ResponseGetConfigurationViewTypeJSON]("success"))
		return
	default:
		resp, err := h.configurationSvc.GetConfigurationViewTypeSchema(r.Context(), request)
		if err != nil {
			response.Err[string](w,
				response.SetErr[string](err.Error()))
			return
		}

		response.Json[configuratordto.ResponseGetConfigurationsViewTypeSchema](w,
			response.SetData[configuratordto.ResponseGetConfigurationsViewTypeSchema](resp),
			response.SetMessage[configuratordto.ResponseGetConfigurationsViewTypeSchema]("success"))
		return
	}
}

// SetConfiguration set new config
// @Summary set new config
// @Description Set new config
// @Param x-clientkey header string true "<Client Key>"
// @Param RequestSetConfiguration body configuratordto.RequestSetConfiguration true "create new field of config"
// @Tags Config
// @Produce json
// @Router /v1/configurator [POST]
func (h *HttpHandle) SetConfiguration(w http.ResponseWriter, r *http.Request) {
	request := configuratordto.RequestSetConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	err = h.configurationSvc.SetConfiguration(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[string](w,
		response.SetMessage[string]("success"))
}

// UpdateConfiguration update new config
// @Summary update new config
// @Description update new config
// @Tags Config
// @Param x-clientkey header string true "<Client Key>"
// @Param RequestUpdateConfiguration body configuratordto.RequestUpdateConfiguration true "update data of config"
// @Produce json
// @Router /v1/configurator [PUT]
func (h *HttpHandle) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {
	request := configuratordto.RequestUpdateConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
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
// @Description update or set configuration
// @Param x-clientkey header string true "<Client Key>"
// @Param RequestSetConfiguration body configuratordto.RequestSetConfiguration true "create new field of config"
// @Tags Config
// @Produce json
// @Router /v1/configurator/upsert [POST]
func (h *HttpHandle) UpsertConfiguration(w http.ResponseWriter, r *http.Request) {
	request := configuratordto.RequestSetConfiguration{
		XClientKey: r.Header.Get("x-clientkey"),
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
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
// @Description delete a config
// @Param x-clientkey header string true "<Client Key>"
// @Param id path string true "The config field identifier it's a UUID."
// @Tags Config
// @Produce json
// @Router /v1/configurator/{id} [DELETE]
func (h *HttpHandle) DeleteConfiguration(w http.ResponseWriter, r *http.Request) {
	request := configuratordto.RequestDeleteConfiguration{
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

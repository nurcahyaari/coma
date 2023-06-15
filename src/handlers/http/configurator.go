package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	configuratordto "github.com/coma/coma/src/domains/configurator/dto"
)

// SetConfiguration set new config
// @Summary set new config
// @Description Set new config
// @Param RequestSetConfiguration body configuratordto.RequestSetConfiguration true "create new field of config"
// @Tags Config
// @Produce json
// @Router /v1/configurator [POST]
func (h *HttpHandle) SetConfiguration(w http.ResponseWriter, r *http.Request) {
	request := configuratordto.RequestSetConfiguration{}

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
// @Produce json
// @Router /v1/configurator [PUT]
func (h *HttpHandle) UpdateConfiguration(w http.ResponseWriter, r *http.Request) {

}

// GetConfiguration get it's config
// @Summary set new config
// @Description Set new config
// @Tags Config
// @Produce json
// @Router /v1/configurator [GET]
func (h *HttpHandle) GetConfiguration(w http.ResponseWriter, r *http.Request) {

}

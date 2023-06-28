package http

import (
	"encoding/json"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	applicationdto "github.com/coma/coma/src/domains/application/dto"
)

// FindApplicationStages get stages
// @Summary get stages
// @Description get stages
// @Param stageName query string true "<Stage name>"
// @Tags Stages
// @Produce json
// @Router /v1/stages [GET]
func (h *HttpHandle) FindApplicationStages(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestFindStage{
		Name: r.FormValue("stageName"),
	}

	resp, err := h.applicationStageSvc.FindStages(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[applicationdto.ResponseStages](w,
		response.SetMessage[applicationdto.ResponseStages]("success"),
		response.SetData[applicationdto.ResponseStages](resp))
}

// CreateApplicationStages set new stages
// @Summary set new stages
// @Description Set new stages
// @Param RequestCreateStage body applicationdto.RequestCreateStage true "create new environtment"
// @Tags Stages
// @Produce json
// @Router /v1/stages [POST]
func (h *HttpHandle) CreateApplicationStages(w http.ResponseWriter, r *http.Request) {
	request := applicationdto.RequestCreateStage{}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.Err[string](w,
			response.SetMessage[string](err.Error()))
		return
	}

	resp, err := h.applicationStageSvc.CreateStage(r.Context(), request)
	if err != nil {
		response.Err[string](w,
			response.SetErr[string](err.Error()))
		return
	}

	response.Json[applicationdto.ResponseStage](w,
		response.SetMessage[applicationdto.ResponseStage]("success"),
		response.SetData[applicationdto.ResponseStage](resp))
}

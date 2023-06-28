package dto

import (
	"github.com/coma/coma/src/domains/application/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RequestCreateStage struct {
	Name string `json:"name"`
}

func (r RequestCreateStage) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.Name, validation.Required))

	return validation.ValidateStruct(&r, validationFieldRules...)
}

func (r RequestCreateStage) NewApplicationStage() model.ApplicationStage {
	uuid := uuid.New()
	return model.ApplicationStage{
		Id:   uuid.String(),
		Name: r.Name,
	}
}

type ResponseStage struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewResponseStage(data model.ApplicationStage) ResponseStage {
	return ResponseStage{
		Id:   data.Id,
		Name: data.Name,
	}
}

type ResponseStages []ResponseStage

func NewResponseStages(datas model.ApplicationStages) ResponseStages {
	stages := make(ResponseStages, 0)
	for _, data := range datas {
		stages = append(stages, NewResponseStage(data))
	}
	return stages
}

type RequestFindStage struct {
	Name string
}

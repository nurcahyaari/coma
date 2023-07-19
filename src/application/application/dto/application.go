package dto

import (
	"errors"

	"github.com/coma/coma/src/domains/entity"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type ApplicationType string

func (a ApplicationType) String() string {
	return string(a)
}

var MapApplicationType = map[ApplicationType]string{
	"service": "service",
	"client":  "client",
}

var (
	ApplicationTypeService ApplicationType = "service"
	ApplicationTypeClient  ApplicationType = "client"
)

func (a ApplicationType) Validate(value any) error {
	if _, ok := MapApplicationType[a]; !ok {
		return errors.New("err: application type is not found")
	}
	return nil
}

type RequestCreateApplication struct {
	StageId string          `json:"stageId"`
	Type    ApplicationType `json:"type"`
	Name    string          `json:"name"`
}

func (r RequestCreateApplication) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.StageId, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.Name, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.Type, validation.Required, validation.By(r.Type.Validate)))

	return validation.ValidateStruct(&r, validationFieldRules...)
}

func (r RequestCreateApplication) NewApplication() entity.Application {
	uuid := uuid.New()
	return entity.Application{
		Id:      uuid.String(),
		StageId: r.StageId,
		Type:    r.Type.String(),
		Name:    r.Name,
	}
}

type ResponseApplication struct {
	Id      string `json:"id"`
	StageId string `json:"stageId"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

func NewResponseApplication(data entity.Application) ResponseApplication {
	return ResponseApplication{
		Id:      data.Id,
		StageId: data.StageId,
		Name:    data.Name,
		Type:    data.Type,
	}
}

type ResponseApplications []ResponseApplication

func NewResponseApplications(datas entity.Applications) ResponseApplications {
	applications := make(ResponseApplications, 0)
	for _, data := range datas {
		applications = append(applications, NewResponseApplication(data))
	}
	return applications
}

type RequestFindApplication struct {
	Id      string
	Name    string
	StageId string
}

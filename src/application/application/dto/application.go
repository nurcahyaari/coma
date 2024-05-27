package dto

import (
	"errors"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
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
	Type ApplicationType `json:"type"`
	Name string          `json:"name"`
}

func (r RequestCreateApplication) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.Name, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.Type, validation.Required, validation.By(r.Type.Validate)))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))

}

func (r RequestCreateApplication) NewApplication() entity.Application {
	uuid := uuid.New()
	return entity.Application{
		Id:   uuid.String(),
		Type: r.Type.String(),
		Name: r.Name,
	}
}

type ResponseApplication struct {
	Id             string                     `json:"id"`
	Type           string                     `json:"type"`
	Name           string                     `json:"name"`
	ApplicationKey ResponseFindApplicationKey `json:"applicationKey"`
}

func (r *ResponseApplication) AttachApplicationKey(applicationKey ResponseFindApplicationKey) {
	r.ApplicationKey = applicationKey
}

func NewResponseApplication(data entity.Application) ResponseApplication {
	return ResponseApplication{
		Id:   data.Id,
		Name: data.Name,
		Type: data.Type,
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
	Id   string
	Name string
}

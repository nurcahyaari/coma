package dto

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
)

type RequestCreateApplicationKey struct {
	ApplicationId string `json:"applicationId"`
}

func (r RequestCreateApplicationKey) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.ApplicationId, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.New(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestCreateApplicationKey) ApplicationKey() entity.ApplicationKey {
	id := uuid.New()
	return entity.ApplicationKey{
		Id:            id.String(),
		ApplicationId: r.ApplicationId,
	}
}

type ResponseCreateApplicationKey struct {
	ApplicationName string `json:"applicationName"`
	Key             string `json:"key"`
}

type RequestInternalFindApplicationKey struct {
	ApplicationId  string
	Key            string
	SkipValidation bool
}

type RequestFindApplicationKey struct {
	ApplicationId   string `json:"applicationId"`
	ApplicationName string `json:"applicationName"`
	Key             string `json:"key"`
}

func (r RequestFindApplicationKey) ValidateKey() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.Key, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.New(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestFindApplicationKey) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.ApplicationId, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.New(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestFindApplicationKey) FilterApplicationKey() entity.FilterApplicationKey {
	return entity.FilterApplicationKey{
		ApplicationId: r.ApplicationId,
		Key:           r.Key,
	}
}

func (r RequestFindApplicationKey) FilterApplication() entity.FilterApplication {
	return entity.FilterApplication{
		Id: r.ApplicationId,
	}
}

type ResponseFindApplicationKey struct {
	Id              string `json:"id"`
	ApplicationId   string `json:"applicationId"`
	ApplicationName string `json:"applicationName"`
	Key             string `json:"key"`
}

func (s *ResponseFindApplicationKey) AttachApplication(application entity.Application) *ResponseFindApplicationKey {
	s.ApplicationName = application.Name
	return s
}

func NewResponseFindApplicationKey(applicationKey entity.ApplicationKey) ResponseFindApplicationKey {
	return ResponseFindApplicationKey{
		Id:            applicationKey.Id,
		ApplicationId: applicationKey.ApplicationId,
		Key:           applicationKey.Key,
	}
}

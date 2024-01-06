package dto

import (
	"net/http"

	internalerror "github.com/coma/coma/internal/x/errors"
	"github.com/coma/coma/src/domain/entity"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RequestCreateApplicationKey struct {
	ApplicationId string `json:"applicationId"`
	StageId       string `json:"stageId"`
}

func (r RequestCreateApplicationKey) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.ApplicationId, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.StageId, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestCreateApplicationKey) ApplicationKey() entity.ApplicationKey {
	id := uuid.New()
	return entity.ApplicationKey{
		Id:            id.String(),
		ApplicationId: r.ApplicationId,
		StageId:       r.StageId,
	}
}

type ResponseCreateApplicationKey struct {
	ApplicationName string `json:"applicationName"`
	StageName       string `json:"stageName"`
	Key             string `json:"key"`
}

type RequestFindApplicationKey struct {
	ApplicationId   string `json:"applicationId"`
	ApplicationName string `json:"applicationName"`
	StageId         string `json:"stageId"`
	StageName       string `json:"stageName"`
	Key             string `json:"key"`
}

func (r RequestFindApplicationKey) ValidateKey() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.Key, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestFindApplicationKey) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.ApplicationId, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.StageId, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestFindApplicationKey) FilterApplicationKey() entity.FilterApplicationKey {
	return entity.FilterApplicationKey{
		ApplicationId: r.ApplicationId,
		StageId:       r.StageId,
		Key:           r.Key,
	}
}

func (r RequestFindApplicationKey) FilterApplication() entity.FilterApplication {
	return entity.FilterApplication{
		Id: r.ApplicationId,
	}
}

func (r RequestFindApplicationKey) FilterApplicationStage() entity.FilterApplicationStage {
	return entity.FilterApplicationStage{
		Id: r.StageId,
	}
}

type ResponseFindApplicationKey struct {
	Id              string `json:"id"`
	ApplicationId   string `json:"applicationId"`
	ApplicationName string `json:"applicationName"`
	StageId         string `json:"stageId"`
	StageName       string `json:"stageName"`
	Key             string `json:"key"`
}

func (s *ResponseFindApplicationKey) AttachApplication(application entity.Application) *ResponseFindApplicationKey {
	s.ApplicationName = application.Name
	return s
}

func (s *ResponseFindApplicationKey) AttachApplicationStage(applicationStage entity.ApplicationStage) *ResponseFindApplicationKey {
	s.StageName = applicationStage.Name
	return s
}

func NewResponseFindApplicationKey(applicationKey entity.ApplicationKey) ResponseFindApplicationKey {
	return ResponseFindApplicationKey{
		Id:            applicationKey.Id,
		ApplicationId: applicationKey.ApplicationId,
		StageId:       applicationKey.StageId,
		Key:           applicationKey.Key,
	}
}

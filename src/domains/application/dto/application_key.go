package dto

import (
	"github.com/coma/coma/src/domains/application/model"
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

	return validation.ValidateStruct(&r, validationFieldRules...)
}

func (r RequestCreateApplicationKey) ApplicationKey() model.ApplicationKey {
	id := uuid.New()
	return model.ApplicationKey{
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
}

func (r RequestFindApplicationKey) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.ApplicationId, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.StageId, validation.Required))

	return validation.ValidateStruct(&r, validationFieldRules...)
}

func (r RequestFindApplicationKey) FilterApplicationKey() model.FilterApplicationKey {
	return model.FilterApplicationKey{
		ApplicationId: r.ApplicationId,
		StageId:       r.StageId,
	}
}

func (r RequestFindApplicationKey) FilterApplication() model.FilterApplication {
	return model.FilterApplication{
		Id: r.ApplicationId,
	}
}

func (r RequestFindApplicationKey) FilterApplicationStage() model.FilterApplicationStage {
	return model.FilterApplicationStage{
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

func (s *ResponseFindApplicationKey) AttachApplication(application model.Application) *ResponseFindApplicationKey {
	s.ApplicationName = application.Name
	return s
}

func (s *ResponseFindApplicationKey) AttachApplicationStage(applicationStage model.ApplicationStage) *ResponseFindApplicationKey {
	s.StageName = applicationStage.Name
	return s
}

func NewResponseFindApplicationKey(applicationKey model.ApplicationKey) ResponseFindApplicationKey {
	return ResponseFindApplicationKey{
		Id:            applicationKey.Id,
		ApplicationId: applicationKey.ApplicationId,
		StageId:       applicationKey.StageId,
		Key:           applicationKey.Key,
	}
}

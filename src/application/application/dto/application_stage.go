package dto

import (
	"net/http"

	internalerror "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/domain/entity"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RequestCreateStage struct {
	Name string `json:"name"`
}

func (r RequestCreateStage) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.Name, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestCreateStage) NewApplicationStage() entity.ApplicationStage {
	uuid := uuid.New()
	return entity.ApplicationStage{
		Id:   uuid.String(),
		Name: r.Name,
	}
}

type ResponseStage struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func NewResponseStage(data entity.ApplicationStage) ResponseStage {
	return ResponseStage{
		Id:   data.Id,
		Name: data.Name,
	}
}

type ResponseStages []ResponseStage

func NewResponseStages(datas entity.ApplicationStages) ResponseStages {
	stages := make(ResponseStages, 0)
	for _, data := range datas {
		stages = append(stages, NewResponseStage(data))
	}
	return stages
}

type RequestFindStage struct {
	Name string
}

package dto

import (
	"encoding/json"
	"net/http"
	"regexp"

	internalerror "github.com/coma/coma/internal/utils/errors"
	"github.com/coma/coma/src/domain/entity"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RequestSetConfiguration struct {
	XClientKey string `json:"-"`
	Field      string `json:"field"`
	Value      any    `json:"value"`
}

func (r RequestSetConfiguration) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.XClientKey, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.Field, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z_]+$"))))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.Value, validation.Required))

	err := validation.ValidateStruct(&r, validationFieldRules...)
	if err == nil {
		return nil
	}

	return internalerror.NewError(err,
		internalerror.SetErrorCode(http.StatusBadRequest),
		internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
}

func (r RequestSetConfiguration) Configuration() entity.Configuration {
	uuid := uuid.New()
	configuration := entity.Configuration{
		Id:        uuid.String(),
		ClientKey: r.XClientKey,
		Field:     r.Field,
		Value:     r.Value,
	}

	return configuration
}

type ResponseSetConfiguration struct {
	Id   string          `json:"id"`
	Data json.RawMessage `json:"data"`
}

package dto

import (
	"net/http"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
	"github.com/nurcahyaari/coma/src/domain/entity"
)

type RequestUpdateConfiguration struct {
	XClientKey string `json:"-"`
	Id         string `json:"id"`
	Field      string `json:"field"`
	Value      any    `json:"value"`
}

func (r RequestUpdateConfiguration) Validate() error {
	validationFieldRules := []*validation.FieldRules{}

	validationFieldRules = append(validationFieldRules, validation.Field(&r.XClientKey, validation.Required))
	validationFieldRules = append(validationFieldRules, validation.Field(&r.Id, validation.Required))
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

func (r RequestUpdateConfiguration) Configuration() entity.Configuration {
	return entity.Configuration{
		Id:        r.Id,
		ClientKey: r.XClientKey,
		Field:     r.Field,
		Value:     r.Value,
	}
}

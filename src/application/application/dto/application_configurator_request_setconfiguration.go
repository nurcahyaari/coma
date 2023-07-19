package dto

import (
	"encoding/json"
	"regexp"

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

	validationFieldRules = append(validationFieldRules, validation.Field(&r.Field, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z_]+$"))))

	return validation.ValidateStruct(&r, validationFieldRules...)
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

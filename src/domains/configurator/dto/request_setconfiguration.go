package dto

import (
	"github.com/coma/coma/src/domains/configurator/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
)

type RequestSetConfiguration struct {
	XClientKey  string `json:"-"`
	ParentField string `json:"parentField"`
	Field       string `json:"field"`
	Value       string `json:"value"`
}

func (r RequestSetConfiguration) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Field, validation.Required),
		validation.Field(&r.Value, validation.Required),
	)
}

func (r RequestSetConfiguration) Configuration() model.Configuration {
	uuid := uuid.New()
	return model.Configuration{
		Id:        uuid.String(),
		ClientKey: r.XClientKey,
		Field:     r.Field,
		Value:     r.Value,
	}
}

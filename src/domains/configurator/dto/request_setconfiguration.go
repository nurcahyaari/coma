package dto

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type RequestSetConfiguration struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

func (r RequestSetConfiguration) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Field, validation.Required),
		validation.Field(&r.Value, validation.Required),
	)
}

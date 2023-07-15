package dto

import (
	"github.com/coma/coma/src/domains/application/model"
)

type RequestUpdateConfiguration struct {
	XClientKey string `json:"-"`
	Id         string `json:"id"`
	Field      string `json:"field"`
	Value      any    `json:"value"`
}

func (r RequestUpdateConfiguration) Configuration() model.Configuration {
	return model.Configuration{
		Id:        r.Id,
		ClientKey: r.XClientKey,
		Field:     r.Field,
		Value:     r.Value,
	}
}

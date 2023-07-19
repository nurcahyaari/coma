package dto

import "github.com/coma/coma/src/domain/entity"

type RequestUpdateConfiguration struct {
	XClientKey string `json:"-"`
	Id         string `json:"id"`
	Field      string `json:"field"`
	Value      any    `json:"value"`
}

func (r RequestUpdateConfiguration) Configuration() entity.Configuration {
	return entity.Configuration{
		Id:        r.Id,
		ClientKey: r.XClientKey,
		Field:     r.Field,
		Value:     r.Value,
	}
}

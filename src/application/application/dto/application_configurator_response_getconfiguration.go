package dto

import (
	"encoding/json"

	"github.com/coma/coma/src/domain/entity"
	"gopkg.in/guregu/null.v4"
)

type ResponseGetConfigurationViewTypeJSON struct {
	ClientKey string          `json:"clientKey"`
	Data      json.RawMessage `json:"data"`
}

func (r *ResponseGetConfigurationViewTypeJSON) SetData(data entity.Configurations) error {
	var (
		mapFieldValue = make(map[string]interface{})
	)

	if len(data) == 0 {
		return nil
	}

	for _, d := range data {
		mapFieldValue[d.Field] = d.Value
	}

	byt, err := json.Marshal(mapFieldValue)
	if err != nil {
		return err
	}

	r.Data = byt
	return nil
}

func NewResponseGetConfigurationViewTypeJSON(clientKey string) ResponseGetConfigurationViewTypeJSON {
	return ResponseGetConfigurationViewTypeJSON{
		ClientKey: clientKey,
	}
}

type ResponseGetConfigurationViewTypeSchema struct {
	Id          string      `json:"id"`
	ClientKey   string      `json:"clientKey"`
	ParentField null.String `json:"parentField"`
	Field       string      `json:"field"`
	Value       any         `json:"value"`
}

func NewResponseGetConfigurationViewTypeSchema(data entity.Configuration) ResponseGetConfigurationViewTypeSchema {
	return ResponseGetConfigurationViewTypeSchema{
		Id:        data.Id,
		ClientKey: data.ClientKey,
		Field:     data.Field,
		Value:     data.Value,
	}
}

type ResponseGetConfigurationsViewTypeSchema []ResponseGetConfigurationViewTypeSchema

func NewResponseGetConfigurationsViewTypeSchema(data entity.Configurations) ResponseGetConfigurationsViewTypeSchema {
	responses := make(ResponseGetConfigurationsViewTypeSchema, 0)
	for _, d := range data {
		responses = append(responses, NewResponseGetConfigurationViewTypeSchema(d))
	}
	return responses
}

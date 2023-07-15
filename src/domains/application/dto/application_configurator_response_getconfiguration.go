package dto

import (
	"encoding/json"

	"github.com/coma/coma/src/domains/application/model"
	"gopkg.in/guregu/null.v4"
)

type ResponseGetConfigurationViewTypeJSON struct {
	ClientKey string          `json:"clientKey"`
	Data      json.RawMessage `json:"data"`
}

type ResponseGetConfiguratorViewType interface {
	model.Configurations
}

func NewResponseGetConfigurationViewTypeJSON[T ResponseGetConfiguratorViewType](data T) (ResponseGetConfigurationViewTypeJSON, error) {
	var (
		response      = ResponseGetConfigurationViewTypeJSON{}
		mapFieldValue = make(map[string]interface{})
	)

	if len(data) == 0 {
		return response, nil
	}

	for _, d := range data {
		mapFieldValue[d.Field] = d.Value
	}

	byt, err := json.Marshal(mapFieldValue)
	if err != nil {
		return response, err
	}

	response.ClientKey = data[0].ClientKey
	response.Data = byt

	return response, nil
}

type ResponseGetConfigurationViewTypeSchema struct {
	Id          string      `json:"id"`
	ClientKey   string      `json:"clientKey"`
	ParentField null.String `json:"parentField"`
	Field       string      `json:"field"`
	Value       any         `json:"value"`
}

func NewResponseGetConfigurationViewTypeSchema(data model.Configuration) ResponseGetConfigurationViewTypeSchema {
	return ResponseGetConfigurationViewTypeSchema{
		Id:        data.Id,
		ClientKey: data.ClientKey,
		Field:     data.Field,
		Value:     data.Value,
	}
}

type ResponseGetConfigurationsViewTypeSchema []ResponseGetConfigurationViewTypeSchema

func NewResponseGetConfigurationsViewTypeSchema(data model.Configurations) ResponseGetConfigurationsViewTypeSchema {
	responses := make(ResponseGetConfigurationsViewTypeSchema, 0)
	for _, d := range data {
		responses = append(responses, NewResponseGetConfigurationViewTypeSchema(d))
	}
	return responses
}
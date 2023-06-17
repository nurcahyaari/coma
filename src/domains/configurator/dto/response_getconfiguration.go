package dto

import (
	"encoding/json"

	"github.com/coma/coma/src/domains/configurator/model"
)

type ResponseGetClientConfiguration struct {
	ClientKey string          `json:"clientKey"`
	Data      json.RawMessage `json:"data"`
}

type ResponseGetClientConfigurator interface {
	model.Configurations
}

func NewResponseGetClientConfiguration[T ResponseGetClientConfigurator](data T) (ResponseGetClientConfiguration, error) {
	var (
		response      = ResponseGetClientConfiguration{}
		mapFieldValue = make(map[string]string)
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

type ResponseGetConfiguration struct {
	Id        string `json:"id"`
	ClientKey string `json:"clientKey"`
	Field     string `json:"field"`
	Value     string `json:"value"`
}

func NewResponseGetConfiguration(data model.Configuration) ResponseGetConfiguration {
	return ResponseGetConfiguration{
		Id:        data.Id,
		ClientKey: data.ClientKey,
		Field:     data.Field,
		Value:     data.Value,
	}
}

type ResponseGetConfigurations []ResponseGetConfiguration

func NewResponseGetConfigurations(data model.Configurations) ResponseGetConfigurations {
	responses := make(ResponseGetConfigurations, 0)
	for _, d := range data {
		responses = append(responses, NewResponseGetConfiguration(d))
	}
	return responses
}

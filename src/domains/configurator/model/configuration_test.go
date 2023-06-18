package model_test

import (
	"testing"

	"github.com/coma/coma/src/domains/configurator/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestConfiguration(t *testing.T) {
	testCases := []struct {
		name     string
		mock     model.Configuration
		expected map[string]interface{}
		withErr  bool
		actual   func(data model.Configuration) (map[string]interface{}, error)
	}{
		{
			name: "test1",
			mock: model.Configuration{
				Id:        "1",
				ClientKey: "1",
				Field:     "1",
				Value:     null.StringFrom("1"),
			},
			expected: map[string]interface{}{
				"id":        "1",
				"clientKey": "1",
				"field":     "1",
				"value":     "1",
			},
			withErr: false,
			actual: func(data model.Configuration) (map[string]interface{}, error) {
				return data.MapStringInterface()
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.actual(test.mock)

			if test.withErr {
				assert.Error(t, err)
			}

			assert.Equal(t, test.expected, actual)

		})
	}
}

func TestUpdateConfiguration(t *testing.T) {
	testCases := []struct {
		name     string
		expected model.Configurations
		actual   func() model.Configurations
	}{
		{
			name: "test update configuration no update",
			expected: model.Configurations{
				{
					Id:    "1",
					Field: "app",
					Value: null.StringFrom("test"),
				},
				{
					Id:    "2",
					Field: "port",
					Value: null.StringFrom("1234"),
				},
			},
			actual: func() model.Configurations {
				oldData := model.Configurations{
					{
						Id:    "1",
						Field: "app",
						Value: null.StringFrom("test"),
					},
					{
						Id:    "2",
						Field: "port",
						Value: null.StringFrom("1234"),
					},
				}
				newData := model.Configurations{}

				oldData.Update(newData.MapConfigurationById())

				return oldData
			},
		},
		{
			name: "test update configuration",
			expected: model.Configurations{
				{
					Id:    "1",
					Field: "app",
					Value: null.StringFrom("test"),
				},
				{
					Id:    "2",
					Field: "port",
					Value: null.StringFrom("1234"),
				},
			},
			actual: func() model.Configurations {
				oldData := model.Configurations{
					{
						Id:    "1",
						Field: "app",
						Value: null.StringFrom("test"),
					},
					{
						Id:    "2",
						Field: "ports",
						Value: null.StringFrom("12312"),
					},
				}
				newData := model.Configurations{
					{
						Id:    "2",
						Field: "port",
						Value: null.StringFrom("1234"),
					},
				}

				oldData.Update(newData.MapConfigurationById())

				return oldData
			},
		},
		{
			name: "test update configuration all data",
			expected: model.Configurations{
				{
					Id:    "1",
					Field: "app",
					Value: null.StringFrom("test"),
				},
				{
					Id:    "2",
					Field: "port",
					Value: null.StringFrom("1234"),
				},
			},
			actual: func() model.Configurations {
				oldData := model.Configurations{
					{
						Id:    "1",
						Field: "app",
						Value: null.StringFrom("test1"),
					},
					{
						Id:    "2",
						Field: "ports",
						Value: null.StringFrom("12312"),
					},
				}
				newData := model.Configurations{
					{
						Id:    "1",
						Field: "app",
						Value: null.StringFrom("test"),
					},
					{
						Id:    "2",
						Field: "port",
						Value: null.StringFrom("1234"),
					},
				}

				oldData.Update(newData.MapConfigurationById())

				return oldData
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := test.actual()

			assert.Equal(t, test.expected, actual)

		})
	}
}

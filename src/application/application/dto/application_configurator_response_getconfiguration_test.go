package dto_test

import (
	"testing"

	"github.com/nurcahyaari/coma/src/application/application/dto"
	"github.com/nurcahyaari/coma/src/domain/entity"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"
)

func TestNewResponseGetConfiguration(t *testing.T) {
	testCases := []struct {
		name      string
		expected  dto.ResponseGetConfigurationViewTypeJSON
		haveError bool
		actual    func() (dto.ResponseGetConfigurationViewTypeJSON, error)
	}{
		{
			name:      "response from entity.Configurations when empty",
			haveError: false,
			expected: dto.ResponseGetConfigurationViewTypeJSON{
				Data: []byte(nil),
			},
			actual: func() (dto.ResponseGetConfigurationViewTypeJSON, error) {
				response := dto.NewResponseGetConfigurationViewTypeJSON("")
				return response, nil
			},
		},
		{
			name:      "response from entity.Configurations",
			haveError: false,
			expected: dto.ResponseGetConfigurationViewTypeJSON{
				ClientKey: "1",
				Data:      []byte(`{"age":"1","name":"test"}`),
			},
			actual: func() (dto.ResponseGetConfigurationViewTypeJSON, error) {
				response := dto.NewResponseGetConfigurationViewTypeJSON("1")
				response.SetData(entity.Configurations{
					{
						Id:        "1",
						ClientKey: "1",
						Field:     "name",
						Value:     null.StringFrom("test"),
					},
					{
						Id:        "2",
						ClientKey: "1",
						Field:     "age",
						Value:     null.StringFrom("1"),
					},
				})
				return response, nil
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual, err := test.actual()

			if test.haveError {
				assert.Error(t, err)
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

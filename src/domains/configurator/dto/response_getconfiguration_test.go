package dto_test

import (
	"testing"

	"github.com/coma/coma/src/domains/configurator/dto"
	"github.com/coma/coma/src/domains/configurator/model"
	"github.com/stretchr/testify/assert"
)

func TestNewResponseGetConfiguration(t *testing.T) {
	testCases := []struct {
		name      string
		expected  dto.ResponseGetClientConfiguration
		haveError bool
		actual    func() (dto.ResponseGetClientConfiguration, error)
	}{
		{
			name:      "response from model.Configurations when empty",
			haveError: false,
			expected: dto.ResponseGetClientConfiguration{
				Data: []byte(nil),
			},
			actual: func() (dto.ResponseGetClientConfiguration, error) {
				return dto.NewResponseGetClientConfiguration[model.Configurations](
					model.Configurations{},
				)
			},
		},
		{
			name:      "response from model.Configurations",
			haveError: false,
			expected: dto.ResponseGetClientConfiguration{
				ClientKey: "1",
				Data:      []byte(`{"age":"1","name":"test"}`),
			},
			actual: func() (dto.ResponseGetClientConfiguration, error) {
				return dto.NewResponseGetClientConfiguration[model.Configurations](
					model.Configurations{
						{
							Id:        "1",
							ClientKey: "1",
							Field:     "name",
							Value:     "test",
						},
						{
							Id:        "2",
							ClientKey: "1",
							Field:     "age",
							Value:     "1",
						},
					},
				)
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

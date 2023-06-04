package dto_test

import (
	"errors"
	"testing"

	"github.com/coma/coma/src/domains/distributor/dto"
	"github.com/stretchr/testify/assert"
)

func TestRequestDistributeValidate(t *testing.T) {
	testCases := []struct {
		name     string
		expected []error
		actual   func() []error
	}{
		{
			name:     "test1-valid all",
			expected: nil,
			actual: func() []error {
				return dto.RequestDistribute{
					ApiToken:    "12345",
					Data:        []byte(`"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  }\n}"`),
					ContentType: "json",
				}.Validate()
			},
		},
		{
			name: "test2 - api token empty",
			expected: []error{
				errors.New("api token cannot be nulled or empty"),
			},
			actual: func() []error {
				return dto.RequestDistribute{
					ApiToken:    "",
					Data:        []byte(`"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  }\n}"`),
					ContentType: "json",
				}.Validate()
			},
		},
		{
			name: "test3 - content type empty",
			expected: []error{
				errors.New("content type cannot be nulled or empty"),
			},
			actual: func() []error {
				return dto.RequestDistribute{
					ApiToken:    "1234",
					Data:        []byte(`"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  }\n}"`),
					ContentType: "",
				}.Validate()
			},
		},
		{
			name: "test4 - data is not valid json",
			expected: []error{
				errors.New("data must be a valid json"),
			},
			actual: func() []error {
				return dto.RequestDistribute{
					ApiToken:    "1234",
					Data:        nil,
					ContentType: "json",
				}.Validate()
			},
		},
		{
			name: "test5 - invalid all",
			expected: []error{
				errors.New("content type cannot be nulled or empty"),
				errors.New("api token cannot be nulled or empty"),
				errors.New("data must be a valid json"),
			},
			actual: func() []error {
				return dto.RequestDistribute{
					ApiToken:    "",
					Data:        nil,
					ContentType: "",
				}.Validate()
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := test.actual()
			expected := test.expected

			assert.Equal(t, expected, actual)
		})
	}
}

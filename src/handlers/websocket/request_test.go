package websocket_test

import (
	"errors"
	"testing"

	"github.com/nurcahyaari/coma/src/handlers/websocket"
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
				return websocket.RequestDistribute{
					ClientKey: "12345",
					Data:      []byte(`"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  }\n}"`),
				}.Validate()
			},
		},
		{
			name: "test2 - api token empty",
			expected: []error{
				errors.New("api token cannot be nulled or empty"),
			},
			actual: func() []error {
				return websocket.RequestDistribute{
					ClientKey: "",
					Data:      []byte(`"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  }\n}"`),
				}.Validate()
			},
		},
		{
			name: "test3 - invalid all",
			expected: []error{
				errors.New("api token cannot be nulled or empty"),
				errors.New("data must be a valid json"),
			},
			actual: func() []error {
				return websocket.RequestDistribute{
					ClientKey: "",
					Data:      []byte(""),
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

package errors_test

import (
	"errors"
	"testing"

	internalerror "github.com/coma/coma/internal/utils/errors"
	"github.com/stretchr/testify/assert"
)

func TestErrorAsObject(t *testing.T) {
	testCases := []struct {
		name     string
		expected any
		actual   func() any
	}{
		{
			name: "test1 - ozzo error",
			expected: map[string]string{
				"Zip":    "cannot be blank",
				"City":   "cannot be blank",
				"Street": "cannot be blank",
				"State":  "cannot be blank",
			},
			actual: func() any {
				add := Address{}
				err := add.Validate()
				err = internalerror.NewError[map[string]string](err,
					internalerror.SetErrorSource[map[string]string](internalerror.OZZO_VALIDATION_ERR))
				errCustom := err.(*internalerror.Error[map[string]string])
				return errCustom.ErrorAsObject()
			},
		},
		{
			name: "test2 - plain error",
			expected: map[string]string{
				"": "Name cannot be blank",
			},
			actual: func() any {
				err := internalerror.NewError[map[string]string](errors.New("Name cannot be blank"),
					internalerror.SetErrorSource[map[string]string](internalerror.PLAIN_ERR))
				errCustom := err.(*internalerror.Error[map[string]string])
				return errCustom.ErrorAsObject()
			},
		},
		{
			name: "test3 - plain error with field",
			expected: map[string]string{
				"Name": "cannot be blank",
			},
			actual: func() any {
				err := internalerror.NewError[map[string]string](errors.New("Name cannot be blank"),
					internalerror.SetErrorSource[map[string]string](internalerror.PLAIN_ERR),
					internalerror.WithField[map[string]string](true))
				errCustom := err.(*internalerror.Error[map[string]string])
				return errCustom.ErrorAsObject()
			},
		},
		{
			name:     "test4 - without source",
			expected: errors.New("Name cannot be blank"),
			actual: func() any {
				err := internalerror.NewError[error](errors.New("Name cannot be blank"))
				errCustom := err.(*internalerror.Error[error])
				return errCustom.ErrorAsObject()
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := test.actual()
			expected := test.expected
			assert.Equal(t, actual, expected)
		})
	}
}

package errors_test

import (
	"errors"
	"testing"

	internalerror "github.com/nurcahyaari/coma/internal/x/errors"
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
				err = internalerror.New(err,
					internalerror.SetErrorSource(internalerror.OZZO_VALIDATION_ERR))
				errCustom := err.(*internalerror.Error)
				return errCustom.ErrorAsObject()
			},
		},
		{
			name: "test2 - plain error",
			expected: map[string]string{
				"": "Name cannot be blank",
			},
			actual: func() any {
				err := internalerror.New(errors.New("Name cannot be blank"),
					internalerror.SetErrorSource(internalerror.PLAIN_ERR_TEXT))
				errCustom := err.(*internalerror.Error)
				return errCustom.ErrorAsObject()
			},
		},
		{
			name: "test3 - plain error with field",
			expected: map[string]string{
				"Name": "cannot be blank",
			},
			actual: func() any {
				err := internalerror.New(errors.New("Name cannot be blank"),
					internalerror.SetErrorSource(internalerror.PLAIN_ERR_TEXT),
					internalerror.WithField(true))
				errCustom := err.(*internalerror.Error)
				return errCustom.ErrorAsObject()
			},
		},
		{
			name:     "test4 - without source",
			expected: errors.New("Name cannot be blank"),
			actual: func() any {
				err := internalerror.New(errors.New("Name cannot be blank"))
				errCustom := err.(*internalerror.Error)
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

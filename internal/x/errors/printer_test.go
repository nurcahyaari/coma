package errors_test

import (
	"errors"
	"regexp"
	"testing"

	internalerror "github.com/coma/coma/internal/x/errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/stretchr/testify/assert"
)

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func (a Address) Validate() error {
	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street, validation.Required, validation.Length(5, 50)),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City, validation.Required, validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}

func TestValidationPrinterWithOzzo(t *testing.T) {
	testCases := []struct {
		name     string
		expected internalerror.Printer
		actual   func() internalerror.Printer
	}{
		{
			name: "from ozzo validation",
			expected: internalerror.Printer{
				Failures: map[string][]string{
					"City":   {"cannot be blank"},
					"State":  {"cannot be blank"},
					"Zip":    {"cannot be blank"},
					"Street": {"cannot be blank"},
				},
			},
			actual: func() internalerror.Printer {
				v := internalerror.NewPrinter()
				add := Address{}
				err := add.Validate()
				v.OzzoValidationErr(err)
				return v
			},
		},
		{
			name: "from plain string",
			expected: internalerror.Printer{
				Failures: map[string][]string{
					"": {"values cannot be null", "name cannot be null"},
				},
			},
			actual: func() internalerror.Printer {
				v := internalerror.NewPrinter()
				v.PlainErr(errors.New("values cannot be null"), false)
				v.PlainErr(errors.New("name cannot be null"), false)
				return v
			},
		},
		{
			name: "from plain string with field",
			expected: internalerror.Printer{
				Failures: map[string][]string{
					"values": {"cannot be null"},
					"name":   {"cannot be null"},
				},
			},
			actual: func() internalerror.Printer {
				v := internalerror.NewPrinter()
				v.PlainErr(errors.New("values cannot be null"), true)
				v.PlainErr(errors.New("name cannot be null"), true)
				return v
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := test.actual()
			expected := test.expected
			assert.Equal(t, actual.Failures, expected.Failures)
		})
	}
}

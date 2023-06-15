package validation_test

import (
	"regexp"
	"testing"

	internalvalidation "github.com/coma/coma/internal/utils/validation"
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
		expected internalvalidation.Printer
		actual   func() internalvalidation.Printer
	}{
		{
			name: "from ozzo validation",
			expected: internalvalidation.Printer{
				Failures: []string{
					"City: cannot be blank.",
					"State: cannot be blank.",
					"Zip: cannot be blank.",
					"Street: cannot be blank.",
				},
			},
			actual: func() internalvalidation.Printer {
				v := internalvalidation.New()
				add := Address{}
				err := add.Validate()
				v.OzzoValidationErr(err)
				return v
			},
		},
		{
			name: "from plain string",
			expected: internalvalidation.Printer{
				Failures: []string{
					"values cannot be null.",
					"name cannot be null.",
				},
			},
			actual: func() internalvalidation.Printer {
				v := internalvalidation.New()
				v.AppendFailure("values cannot be null")
				v.AppendFailure("name cannot be null")
				return v
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			actual := test.actual()
			expected := test.expected
			assert.ElementsMatch(t, actual.Failures, expected.Failures)
		})
	}
}

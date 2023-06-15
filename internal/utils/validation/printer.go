package validation

import (
	"errors"
	"fmt"
	"strings"
)

type Printer struct {
	Err      error
	Failures []string
}

func New() Printer {
	return Printer{
		Err: errors.New("error: validation"),
	}
}

func (v *Printer) OzzoValidationErr(err error) {
	msg := err.Error()
	msgs := strings.Split(msg, ";")

	for _, msg := range msgs {
		v.AppendFailure(
			strings.TrimSuffix(
				strings.TrimSuffix(strings.TrimPrefix(msg, " "), " "), "."),
		)
	}
}

func (v *Printer) PlainErr(err error) {
	msg := err.Error()
	v.AppendFailure(msg)
}

func (v *Printer) AppendFailure(message string) {
	v.Failures = append(v.Failures, fmt.Sprintf("%s.", message))
}

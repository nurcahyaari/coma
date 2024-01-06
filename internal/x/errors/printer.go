package errors

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Printer struct {
	Err      error
	Failures map[string][]string
}

func NewPrinter() Printer {
	return Printer{
		Err:      errors.New("error: validation"),
		Failures: make(map[string][]string),
	}
}

func (v *Printer) OzzoValidationErr(err error) {
	for field, fieldErr := range err.(validation.Errors) {
		message := fieldErr.Error()
		v.AppendFailure(field, message)
	}
}

func (v *Printer) PlainErr(err error, withField bool) {
	var (
		msg   = err.Error()
		field string
	)

	if withField {
		msgSplit := strings.Split(msg, ";")
		if len(msgSplit) == 1 {
			msgSplit = strings.SplitN(msg, " ", 2)
		}

		field = msg

		if len(msgSplit) > 1 {
			field = msgSplit[0]
			msg = msgSplit[1]
		}
	}

	v.AppendFailure(field, msg)
}

func (v *Printer) AppendFailure(field, message string) {
	v.Failures[field] = append(v.Failures[field], message)
}

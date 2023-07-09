package errors

import (
	"strings"
)

type ErrorSource string

const (
	PLAIN_ERR_TEXT      ErrorSource = "PLAIN_ERR_TEXT"
	OZZO_VALIDATION_ERR ErrorSource = "OZZO_VALIDATION_ERR"
)

type Error struct {
	ErrorSource ErrorSource
	WithField   bool
	Err         error
	printer     Printer
}

func (r *Error) Error() string {
	return r.Err.Error()
}

func (r *Error) ErrorAsObject() any {
	mapString := make(map[string]string)

	switch r.ErrorSource {
	case OZZO_VALIDATION_ERR:
		r.printer.OzzoValidationErr(r.Err)
	case PLAIN_ERR_TEXT:
		r.printer.PlainErr(r.Err, r.WithField)
	default:
		return r.Err
	}

	for field, failure := range r.printer.Failures {
		mapString[field] = strings.Join(failure, ";")
	}

	return mapString
}

type ErrorOpt func(err *Error)

func WithField(withField bool) ErrorOpt {
	return func(err *Error) {
		err.WithField = withField
	}
}

func SetErrorSource(errSource ErrorSource) ErrorOpt {
	return func(err *Error) {
		err.ErrorSource = errSource
	}
}

func NewError(err error, opts ...ErrorOpt) error {
	resp := &Error{
		WithField: false,
		Err:       err,
		printer:   NewPrinter(),
	}

	for _, opt := range opts {
		opt(resp)
	}

	return resp
}

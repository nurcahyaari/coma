package errors

import (
	"strings"
)

type ErrorSource string

const (
	PLAIN_ERR           ErrorSource = "PLAIN_ERR"
	OZZO_VALIDATION_ERR ErrorSource = "OZZO_VALIDATION_ERR"
)

type Error[T any] struct {
	ErrorSource ErrorSource
	WithField   bool
	Err         error
	printer     Printer
}

func (r *Error[T]) Error() string {
	return r.Err.Error()
}

func (r *Error[T]) ErrorAsObject() any {
	mapString := make(map[string]string)

	switch r.ErrorSource {
	case OZZO_VALIDATION_ERR:
		r.printer.OzzoValidationErr(r.Err)
	case PLAIN_ERR:
		r.printer.PlainErr(r.Err, r.WithField)
	default:
		return r.Err
	}

	for field, failure := range r.printer.Failures {
		mapString[field] = strings.Join(failure, ";")
	}

	return mapString
}

type ErrorOpt[T any] func(err *Error[T])

func WithField[T any](withField bool) ErrorOpt[T] {
	return func(err *Error[T]) {
		err.WithField = withField
	}
}

func SetErrorSource[T any](errSource ErrorSource) ErrorOpt[T] {
	return func(err *Error[T]) {
		err.ErrorSource = errSource
	}
}

func NewError[T any](err error, opts ...ErrorOpt[T]) error {
	resp := &Error[T]{
		WithField: false,
		Err:       err,
		printer:   New(),
	}

	for _, opt := range opts {
		opt(resp)
	}

	return resp
}

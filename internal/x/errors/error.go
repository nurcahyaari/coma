package errors

import (
	"net/http"
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
	ErrCode     int
	printer     Printer
}

func (r *Error) Error() string {
	if r.Err == nil {
		return ""
	}
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
		return r.Err.Error()
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

func SetErrorCode(code int) ErrorOpt {
	return func(err *Error) {
		err.ErrCode = code
	}
}

// New will create an error object
// it also will print the stack trace
func New(err error, opts ...ErrorOpt) error {
	if err == nil {
		return nil
	}

	resp := &Error{
		WithField: false,
		Err:       err,
		ErrCode:   http.StatusInternalServerError,
		printer:   NewPrinter(),
	}

	for _, opt := range opts {
		opt(resp)
	}

	return resp
}

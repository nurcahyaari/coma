package errors

import (
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/ztrue/tracerr"
)

type ErrorSource string

const (
	PLAIN_ERR_TEXT      ErrorSource = "PLAIN_ERR_TEXT"
	OZZO_VALIDATION_ERR ErrorSource = "OZZO_VALIDATION_ERR"
)

type Error struct {
	ErrorSource ErrorSource
	WithField   bool
	WithTracer  bool
	Err         error
	ErrCode     int
	printer     Printer
	tracer      tracerr.Error
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

func WithTracer(WithTracer bool) ErrorOpt {
	return func(err *Error) {
		err.WithTracer = WithTracer
	}
}

// New will create an error object
// it also will print the stack trace
func New(err error, opts ...ErrorOpt) error {
	if err == nil {
		return nil
	}

	tracerWrapper := tracerr.Wrap(err)
	frames := tracerr.StackTrace(tracerWrapper)
	tracerErr := tracerr.CustomError(tracerWrapper, frames[:3])

	resp := &Error{
		WithField: false,
		Err:       err,
		ErrCode:   http.StatusInternalServerError,
		printer:   NewPrinter(),
		tracer:    tracerErr,
	}

	for _, opt := range opts {
		opt(resp)
	}

	if resp.WithTracer {
		log.Info().Msg(tracerr.SprintSource(resp.tracer, 1, 1))
	}

	return resp
}

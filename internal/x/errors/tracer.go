package errors

import (
	"github.com/rs/zerolog/log"
	"github.com/ztrue/tracerr"
)

type Trace struct {
	LineAfter   int
	LineBefore  int
	TotalFrames int
}

type TraceOpt func(t *Trace)

// Define total line after the error trace
func LineAfter(n int) TraceOpt {
	return func(t *Trace) {
		t.LineAfter = n
	}
}

func LineBefore(n int) TraceOpt {
	return func(t *Trace) {
		t.LineBefore = n
	}
}

func TotalFrames(n int) TraceOpt {
	return func(t *Trace) {
		t.TotalFrames = n
	}
}

func StackTrace(err error, opts ...TraceOpt) {
	trace := Trace{
		LineAfter:   1,
		LineBefore:  1,
		TotalFrames: 3,
	}

	for _, opt := range opts {
		opt(&trace)
	}

	tracerWrapper := tracerr.Wrap(err)
	frames := tracerr.StackTrace(tracerWrapper)
	tracerErr := tracerr.CustomError(tracerWrapper, frames[:trace.TotalFrames])
	log.Info().Msg(tracerr.SprintSource(tracerErr, trace.LineBefore, trace.LineAfter))
}

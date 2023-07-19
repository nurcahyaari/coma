package pubsub

import (
	"io"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type SubscriberHandler func(id string, r io.Reader)

type SubscriberOpt func(s *subscriber)

func SetSubscriberMaxWorker(maxWorker int) SubscriberOpt {
	return func(s *subscriber) {
		s.maxWorker = maxWorker
	}
}

func SetSubscriberAsync(async bool) SubscriberOpt {
	return func(s *subscriber) {
		s.async = async
	}
}

func SetSubscriberMaxElapsedTime(maxElapsedTime time.Duration) SubscriberOpt {
	return func(s *subscriber) {
		s.maxElapsedTime = maxElapsedTime
	}
}

func SetSubscriberRetryWaitTime(retryWaitTime time.Duration) SubscriberOpt {
	return func(s *subscriber) {
		s.retryWaitTime = retryWaitTime
	}
}

type subscriber struct {
	shutdownDispatcher chan bool
	shutdownConsumer   chan bool
	async              bool
	maxWorker          int
	maxElapsedTime     time.Duration
	retryWaitTime      time.Duration
	handler            SubscriberHandler
	message            chan io.Reader
}

func newSubscriber() *subscriber {
	sub := &subscriber{
		shutdownDispatcher: make(chan bool),
		shutdownConsumer:   make(chan bool),
		async:              false,
		maxWorker:          1,
		maxElapsedTime:     1 * time.Second,
		retryWaitTime:      3 * time.Second,
		message:            make(chan io.Reader),
	}

	if sub.maxWorker == 0 {
		sub.maxWorker = 1
	}

	return sub
}

type SubscriberOption func(s *subscriber)

func PubsubSetMaxWorker(max int) SubscriberOption {
	return func(ps *subscriber) {
		ps.maxWorker = max
	}
}

func PubsubSetAsyncProcess(ok bool) SubscriberOption {
	return func(ps *subscriber) {
		ps.async = ok
	}
}

func PubsubSetMaxElapsedTime(duration time.Duration) SubscriberOption {
	return func(ps *subscriber) {
		ps.maxElapsedTime = duration
	}
}

func PubsubSetRetryWaitTime(duration time.Duration) SubscriberOption {
	return func(ps *subscriber) {
		ps.retryWaitTime = duration
	}
}

func (s *subscriber) registerSubscriberHandler(handler SubscriberHandler, opts ...SubscriberOption) {
	for _, opt := range opts {
		opt(s)
	}
	s.handler = handler
}

func (s *subscriber) shutdownSubscriber() {
	s.shutdownDispatcher <- true
	s.shutdownConsumer <- true
	close(s.shutdownConsumer)
	close(s.shutdownDispatcher)
	close(s.message)
}

func (s *subscriber) listen() {
	for i := 1; i <= s.maxWorker; i++ {
		go s.consume()
	}
}

func (s *subscriber) dispatcher(publisher *publisher) {
	for {
		select {
		case <-s.shutdownDispatcher:
			return
		case message := <-publisher.message:
			s.message <- message
		}
	}
}

func (s *subscriber) consume() {
	for {
		select {
		case <-s.shutdownConsumer:
			return
		case message := <-s.message:
			id := uuid.New()

			backoffExponential := backoff.NewExponentialBackOff()
			backoffExponential.MaxInterval = s.retryWaitTime
			backoffExponential.MaxElapsedTime = s.maxElapsedTime

			if s.async {
				go backoff.Retry(func() error {
					s.handler(id.String(), message)
					return nil
				}, backoffExponential)
				continue
			}

			err := backoff.Retry(func() error {
				s.handler(id.String(), message)
				return nil
			}, backoffExponential)
			if err != nil {
				log.Error().Err(err).Msg("error consume message")
			}
		}
	}
}

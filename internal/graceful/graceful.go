package graceful

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type ServerState int

const (
	StateReady ServerState = iota
	StateShutdown
)

type Operation func(ctx context.Context) error

type RequestGraceful struct {
	WarnPeriod     time.Duration
	ShutdownPeriod time.Duration
	Operations     map[string]Operation
}

func GracefulShutdown(ctx context.Context, req RequestGraceful) {
	if len(req.Operations) == 0 {
		return
	}

	signalchan := make(chan os.Signal, 1)
	signal.Notify(signalchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	oscall := <-signalchan

	wg := sync.WaitGroup{}
	wg.Add(len(req.Operations))
	for k, op := range req.Operations {
		go func(k string, op Operation) {
			defer wg.Done()
			log.Warn().Msgf("Graceful period %s", k)
			time.Sleep(req.WarnPeriod)
			err := op(ctx)
			if err != nil {
				log.Err(err).Msg("Error when stop server")
			}
			log.Warn().Msgf("Shutdown %s", k)
		}(k, op)
	}
	wg.Wait()

	log.Warn().Msg("System was shutdown")
	log.Warn().Msgf("system call:%+v", oscall)
	os.Exit(0)
}

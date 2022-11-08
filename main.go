package main

import (
	"context"
	"sync"

	"github.com/SphericalKat/livechart-go/internal/config"
	"github.com/SphericalKat/livechart-go/internal/lifecycle"
	"github.com/rs/zerolog/log"
)

func main() {
	// load config
	config.Load()

	// create a waitgroup for all tasks
	wg := sync.WaitGroup{}

	// create context for background tasks
	_, cancelFunc := context.WithCancel(context.Background())

	// listen for shutdown signals
	wg.Add(1)
	go lifecycle.ShutdownListener(&wg, &cancelFunc)

	// wait for all tasks to finish
	wg.Wait()

	log.Info().Msg("Graceful shutdown complete")
}

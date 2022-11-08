package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/SphericalKat/livechart-go/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func StartAPI(ctx context.Context, wg *sync.WaitGroup) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	go func() {
		e.Start(fmt.Sprintf(":%s", config.Conf.Port))
	}()

	// listen for context cancellation
	<-ctx.Done()

	// shut down http server
	log.Info().Msg("gracefully shutting down http server...")
	if err := e.Shutdown(context.Background()); err != nil {
		log.Err(err).Msg("server shutdown Failed")
	}

	wg.Done()
}

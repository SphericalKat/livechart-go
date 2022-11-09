package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/SphericalKat/livechart-go/internal/config"
	"github.com/SphericalKat/livechart-go/pkg/controllers"
	"github.com/SphericalKat/livechart-go/pkg/entities"
	"github.com/allegro/bigcache/v3"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
)

var cache *bigcache.BigCache

func init() {
	cache, _ = bigcache.New(context.Background(), bigcache.DefaultConfig(1*time.Hour))
}

func configureEcho() *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.Debug = config.Conf.Env == "dev"

	// set up middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	return e
}

func StartAPI(ctx context.Context, wg *sync.WaitGroup) {
	e := configureEcho()

	e.GET("/", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	e.GET("/latest", func(c echo.Context) error {
		var shows []entities.Show
		cached, err := cache.Get("latest")
		if err != nil {
			shows = controllers.GetLatest()
			ins := []byte{}
			buf := bytes.NewBuffer(ins)
			json.NewEncoder(buf).Encode(&shows)
			cache.Set("latest", buf.Bytes())
		} else {
			json.NewDecoder(bytes.NewReader(cached)).Decode(&shows)
		}

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		c.Response().WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(c.Response())
		encoder.SetEscapeHTML(false)

		if config.Conf.Env == "dev" {
			encoder.SetIndent("", "  ")
		}

		return encoder.Encode(echo.Map{
			"success": true,
			"data":    shows,
		})
	})

	go func() {
		log.Info().Str("addr", fmt.Sprintf("http://localhost:%s", config.Conf.Port)).Msg("started api server")
		if err := e.Start(fmt.Sprintf(":%s", config.Conf.Port)); err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
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

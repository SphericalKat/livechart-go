package config

import (
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Port       string `koanf:"PORT"`
	Env        string `koanf:"ENV"`
	WebsiteURL string `koanf:"WEBSITE_URL"`
}

var Conf *Config

var k = koanf.New(".")

func Load() {
	// set pretty logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// set defaults
	k.Load(confmap.Provider(map[string]interface{}{
		"PORT": 3000,
		"ENV":  "dev",
		"WEBSITE_URL": "https://livechart.me",
	}, "."), nil)

	// attempt to load .env file
	if err := k.Load(file.Provider(".env"), dotenv.Parser()); err != nil {
		log.Info().Err(err).Msg("unable to find env file:")
		log.Info().Msg("falling back to env variables")
	}

	// merge existing env variables
	if err := k.Load(env.Provider("", ".", nil), nil); err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	Conf = &Config{}
	err := k.Unmarshal("", Conf)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// enable debug logging if in dev
	if Conf.Env == "dev" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Msg("loaded config from environment.")
}

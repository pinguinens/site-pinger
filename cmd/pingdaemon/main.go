package main

import (
	"flag"
	"github.com/pinguinens/site-pinger/internal/dialer"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "Custom config path")
	flag.Parse()
}

func main() {
	appConf, err := config.New(configPath)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	appLogger, err := logger.New(appConf.GetLogFile())
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer logger.CloseLogFile()

	appDialer := dialer.New(appConf.DialerTimeout, appConf.DialerTimeout)

	appLogger.Info().Int("status_code", resp.StatusCode).Msg(string(body))
	appLogger.Debug().Msg("finish")
}

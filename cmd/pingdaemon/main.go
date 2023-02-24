package main

import (
	"flag"
	"io"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/daemon"
	"github.com/pinguinens/site-pinger/internal/dialer"
	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/site"
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

	sites, err := site.ParseDir(appConf.SiteDir)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	// TODO: hosts
	hosts := dialer.HostTable{}
	hosts.Add(appConf.Domain, appConf.Hosts[0])

	_ = dialer.New(appConf.DialerTimeout, appConf.DialerTimeout, &hosts)

	response, err := dialer.Ping(appConf.URI)
	if err != nil {
		log.Error().Msg(err.Error())
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Print(err)
	}

	appLogger.Info().Int("status_code", response.StatusCode).Msg(string(body))
	appLogger.Debug().Msg("finish")

	app := daemon.New(appLogger, sites)
	app.Start()
}

package main

import (
	"flag"
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

	appDialer := dialer.New(appConf.DialerTimeout, appConf.DialerKeepAlive)

	app := daemon.New(appLogger, sites, appDialer)
	app.Start()
}

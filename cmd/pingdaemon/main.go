package main

import (
	"flag"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/connector"
	"github.com/pinguinens/site-pinger/internal/daemon"
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
	defer appLogger.Close()

	sites, err := site.ParseDir(appConf.SiteDir)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	hostTables, err := sites.GetHostsList()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	clients := make([]*http.Client, 0, len(hostTables))
	for _, hosts := range hostTables {
		clients = append(clients, connector.New(hosts))
	}

	app := daemon.New(appLogger, clients, sites)
	app.Start()
}

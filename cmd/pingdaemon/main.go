package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/connector"
	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/processor"
	"github.com/pinguinens/site-pinger/internal/service"
	"github.com/pinguinens/site-pinger/internal/site"
)

const (
	appVersion = "0.1.4"
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

	appLogger, err := logger.New(appConf.LogFileName, appConf.ConsoleFormat)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer appLogger.Close()

	appProcessor := processor.New(&appLogger.Logger)

	sites, err := site.ParseDir(appConf.SiteDir)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	hostTables, err := sites.GetHostsList()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	clients := make([]*connector.Connector, 0, len(hostTables))
	for _, hosts := range hostTables {
		clients = append(clients, connector.New(hosts, appVersion))
	}

	app := service.New(appLogger, &appProcessor, clients, sites)
	app.Start()
}

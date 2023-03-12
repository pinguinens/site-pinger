package main

import (
	"flag"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/connector"
	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/messenger"
	dummyMsgSvc "github.com/pinguinens/site-pinger/internal/messenger/dummy"
	"github.com/pinguinens/site-pinger/internal/messenger/msgsvc"
	"github.com/pinguinens/site-pinger/internal/processor"
	"github.com/pinguinens/site-pinger/internal/site"
)

const (
	appVersion = "0.3"
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

	var messegeSrv messenger.Messenger
	if appConf.Messenger.Enabled {
		messegeSrv, err = msgsvc.New(appConf.Messenger.Address, appConf.Messenger.AlarmCodes)
		if err != nil {
			log.Fatal().Msg(err.Error())
		}
		defer messegeSrv.Close()
	} else {
		messegeSrv, _ = dummyMsgSvc.New()
	}

	appProcessor := processor.New(&appLogger.Logger, messegeSrv)

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

	runApp(appLogger, appProcessor, clients, sites)
}

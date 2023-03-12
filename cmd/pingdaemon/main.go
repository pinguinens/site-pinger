package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/connector"
	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/messenger"
	dummyMsgSvc "github.com/pinguinens/site-pinger/internal/messenger/dummy"
	"github.com/pinguinens/site-pinger/internal/messenger/msgsvc"
	"github.com/pinguinens/site-pinger/internal/processor"
	"github.com/pinguinens/site-pinger/internal/service"
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

	app := service.New(appLogger, &appProcessor, clients, sites)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
		<-exit
		cancel()
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Stopping service")
				return
			default:
				app.Start()
				time.Sleep(5 * time.Second)
			}

		}
	}()

	wg.Wait()
	log.Info().Msg("Service stopped")
}

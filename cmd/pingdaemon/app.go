package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/connector"
	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/processor"
	"github.com/pinguinens/site-pinger/internal/service"
	"github.com/pinguinens/site-pinger/internal/site"
)

func runApp(logger *logger.Logger, processor *processor.Processor, clients []*connector.Connector, sites site.Collection, timeout *time.Duration) {
	app := service.New(logger, processor, clients, sites)

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
				log.Info().Msg("Stopping app")
				return
			default:
				app.Start(ctx)
				time.Sleep(*timeout)
			}

		}
	}()

	wg.Wait()
	log.Info().Msg("App stopped")
}

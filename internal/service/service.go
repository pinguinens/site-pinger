package service

import (
	"context"
	"net/url"

	"github.com/pinguinens/site-pinger/internal/connector"
	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/processor"
	"github.com/pinguinens/site-pinger/internal/resource"
	"github.com/pinguinens/site-pinger/internal/site"
)

type Service struct {
	logger    *logger.Logger
	processor *processor.Processor
	clients   []*connector.Connector
	resources []resource.Resource
}

func New(logger *logger.Logger, processor *processor.Processor, clients []*connector.Connector, sites site.Collection) Service {
	resources := make([]resource.Resource, 0, len(clients)*len(sites.List))
	for _, s := range sites.List {
		uri, err := url.Parse(s.Target.URI)
		if err != nil {
			logger.Error().Msg(err.Error())
		}

		for i, h := range s.Target.Hosts {
			resources = append(resources, resource.New(
				s.Target.Method,
				s.Target.URI,
				resource.Host{
					Domain: uri.Host,
					Addr:   h,
				},
				clients[i]),
			)
		}
	}

	return Service{
		logger:    logger,
		processor: processor,
		resources: resources,
		clients:   clients,
	}
}

func (d *Service) Start(ctx context.Context) {
	for _, s := range d.resources {
		select {
		case <-ctx.Done():
			d.logger.Info().Msg("Stopping service...")
			return
		default:
			response, err := s.Ping()
			if err != nil {
				err = d.processor.ProcessError(err)
				if err != nil {
					d.logger.Error().Msg(err.Error())
					continue
				}
				continue
			}

			if response != nil {
				err = d.processor.ProcessResponse(response)
				if err != nil {
					d.logger.Error().Msg(err.Error())
					continue
				}
				err = response.Body.Close()
				if err != nil {
					d.logger.Error().Msg(err.Error())
					continue
				}
			}
		}
	}
}

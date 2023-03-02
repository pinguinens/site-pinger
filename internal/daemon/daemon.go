package daemon

import (
	"io"
	"net/http"
	"net/url"

	"github.com/pinguinens/site-pinger/internal/logger"
	"github.com/pinguinens/site-pinger/internal/resource"
	"github.com/pinguinens/site-pinger/internal/site"
)

type Daemon struct {
	logger    *logger.Logger
	clients   []*http.Client
	resources []resource.Resource
}

func New(logger *logger.Logger, clients []*http.Client, sites site.Collection) Daemon {
	var resources []resource.Resource

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

	return Daemon{
		logger:    logger,
		resources: resources,
		clients:   clients,
	}
}

func (d *Daemon) Start() {
	for _, s := range d.resources {
		response, err := s.Ping()
		if err != nil {
			d.logger.Error().Msg(err.Error())
		}
		if response != nil {
			body, err := io.ReadAll(response.Body)
			if err != nil {
				d.logger.Print(err)
			}

			d.logger.Info().Int("status_code", response.StatusCode).Msg(string(body))
		}
	}
}

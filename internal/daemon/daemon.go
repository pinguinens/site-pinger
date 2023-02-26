package daemon

import (
	"io"
	"net"
	"net/url"

	log "github.com/rs/zerolog"

	"github.com/pinguinens/site-pinger/internal/resource"
	"github.com/pinguinens/site-pinger/internal/site"
)

type Daemon struct {
	logger    *log.Logger
	resources []resource.Resource
}

func New(logger *log.Logger, sites []site.Site, dialer *net.Dialer) Daemon {
	var resources []resource.Resource

	for _, s := range sites {
		uri, err := url.Parse(s.Target.URI)
		if err != nil {
			logger.Error().Msg(err.Error())
		}

		for _, h := range s.Target.Hosts {
			resources = append(resources, resource.New(
				s.Target.Method,
				s.Target.URI,
				resource.Host{
					Domain: uri.Host,
					Addr:   h,
				},
				dialer),
			)
		}
	}

	return Daemon{
		logger:    logger,
		resources: resources,
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

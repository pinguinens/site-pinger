package daemon

import (
	"io"
	"net"
	"net/url"
	"time"

	log "github.com/rs/zerolog"

	"github.com/pinguinens/site-pinger/internal/dialer"
	"github.com/pinguinens/site-pinger/internal/site"
)

type Daemon struct {
	logger *log.Logger
	dialer *net.Dialer
	sites  []site.Site
}

func New(logger *log.Logger, sites []site.Site, dialerTimeout, dialerKeepAlive time.Duration) Daemon {
	hosts := dialer.HostTable{}
	for _, s := range sites {
		uri, err := url.Parse(s.Target.URI)
		if err != nil {
			logger.Error().Msg(err.Error())
		}

		hosts.Add(uri.Host, s.Target.Hosts[0])
	}

	d := dialer.New(dialerTimeout, dialerKeepAlive, &hosts)

	return Daemon{
		logger: logger,
		dialer: d,
		sites:  sites,
	}
}

func (d *Daemon) Start() {
	for _, s := range d.sites {
		response, err := dialer.Ping(s.Target.URI)
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

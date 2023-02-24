package daemon

import (
	log "github.com/rs/zerolog"

	"github.com/pinguinens/site-pinger/internal/site"
)

type Daemon struct {
	logger *log.Logger
	sites  []site.Site
}

func New(logger *log.Logger, sites []site.Site) Daemon {
	return Daemon{
		logger: logger,
		sites:  sites,
	}
}

func (d *Daemon) Start() {}

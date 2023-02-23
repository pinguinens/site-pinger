package daemon

import (
	log "github.com/rs/zerolog"
)

type Daemon struct {
	logger *log.Logger
}

func New(logger *log.Logger) Daemon {
	return Daemon{
		logger: logger,
	}
}

func (d *Daemon) Start() {}

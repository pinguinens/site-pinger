package processor

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	log "github.com/rs/zerolog"

	"github.com/pinguinens/site-pinger/internal/messenger"
)

const (
	AddrField       = "addr"
	MethodField     = "method"
	StatusCodeField = "status_code"
	UrlField        = "url"

	DefaultStatusCode = 0
)

type Processor struct {
	logger    *log.Logger
	messenger *messenger.Messenger
}

func New(logger *log.Logger, messegeSrv *messenger.Messenger) Processor {
	return Processor{
		logger:    logger,
		messenger: messegeSrv,
	}
}

func (p *Processor) ProcessResponse(response *http.Response) error {
	p.logger.Log().Int(StatusCodeField, response.StatusCode).Str(MethodField, response.Request.Method).Str(UrlField, response.Request.URL.String()).Str(AddrField, response.Request.RemoteAddr).Send()

	err := p.messenger.Alarm(response.StatusCode, response.Request.Method, response.Request.URL.String(), response.Request.RemoteAddr)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) ProcessError(err error) error {
	switch e := err.(type) {
	case *url.Error:
		if ie, ok := e.Err.(*net.OpError); ok {
			p.logger.Log().Int(StatusCodeField, DefaultStatusCode).Str(MethodField, strings.ToUpper(e.Op)).Str(UrlField, e.URL).Str(AddrField, ie.Addr.String()).Send()

			err := p.messenger.Alarm(DefaultStatusCode, strings.ToUpper(e.Op), e.URL, ie.Addr.String())
			if err != nil {
				return err
			}

			return nil
		}

		p.logger.Log().Int(StatusCodeField, DefaultStatusCode).Str(MethodField, strings.ToUpper(e.Op)).Str(UrlField, e.URL).Send()
		err := p.messenger.Alarm(DefaultStatusCode, strings.ToUpper(e.Op), e.URL, "")
		if err != nil {
			return err
		}

		return nil
	}

	return err
}

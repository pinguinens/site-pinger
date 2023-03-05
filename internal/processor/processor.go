package processor

import (
	"net"
	"net/http"
	"net/url"
	"strings"

	log "github.com/rs/zerolog"
)

const (
	AddrField       = "addr"
	MethodField     = "method"
	StatusCodeField = "status_code"
	UrlField        = "url"

	DefaultStatusCode = 0

	TextMime = "text"
)

type Processor struct {
	logger *log.Logger
}

func New(logger *log.Logger) Processor {
	return Processor{logger: logger}
}

func (p *Processor) ProcessResponse(response *http.Response) error {
	p.logger.Log().Int(StatusCodeField, response.StatusCode).Str(MethodField, response.Request.Method).Str(UrlField, response.Request.URL.String()).Str(AddrField, response.Request.RemoteAddr).Send()

	return nil
}

func (p *Processor) ProcessError(err error) error {
	switch e := err.(type) {
	case *url.Error:
		if ie, ok := e.Err.(*net.OpError); ok {
			p.logger.Log().Int(StatusCodeField, DefaultStatusCode).Str(MethodField, strings.ToUpper(e.Op)).Str(UrlField, e.URL).Str(AddrField, ie.Addr.String()).Send()

			return nil
		}

		p.logger.Log().Int(StatusCodeField, DefaultStatusCode).Str(MethodField, strings.ToUpper(e.Op)).Str(UrlField, e.URL).Send()

		return nil
	}

	return err
}

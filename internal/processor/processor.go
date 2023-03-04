package processor

import (
	"io"
	"net/http"

	log "github.com/rs/zerolog"
)

const (
	BodyField       = "body"
	MethodField     = "method"
	StatusCodeField = "status_code"
	UrlField        = "url"
)

type Processor struct {
	logger *log.Logger
}

func New(logger *log.Logger) Processor {
	return Processor{logger: logger}
}

func (p *Processor) ProcessResponse(response *http.Response) error {
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	p.logger.Log().Str(MethodField, response.Request.Method).Str(UrlField, response.Request.URL.String()).Int(StatusCodeField, response.StatusCode).Bytes(BodyField, body).Send()

	return nil
}

package processor

import (
	"io"
	"net/http"

	log "github.com/rs/zerolog"
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

	p.logger.Info().Int("status_code", response.StatusCode).Msg(string(body))

	return nil
}

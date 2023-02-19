package main

import (
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/logger"
)

func main() {
	file, err := os.OpenFile("out.log", os.O_WRONLY, 0755)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	logger := logger.New(file)

	resp, err := http.Get("https://example.ru/")
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	logger.Info().Int("status_code", resp.StatusCode).Msg(string(body))
	logger.Debug().Msg("finish")
}

package main

import (
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/logger"
)

type Config struct {
	URL string `yaml:"url"`
}

func main() {
	configBytes, err := os.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("out.log", os.O_WRONLY, 0755)
	if err != nil {
		log.Print(err)
	}
	defer file.Close()

	logger := logger.New(file)

	resp, err := http.Get(config.URL)
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

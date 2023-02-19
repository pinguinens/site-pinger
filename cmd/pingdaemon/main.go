package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/logger"
)

type Config struct {
	URI    string   `yaml:"uri"`
	Domain string   `yaml:"domain"`
	Port   int      `yaml:"port"`
	Hosts  []string `yaml:"hosts"`
}

func init() {
	net.DefaultResolver.PreferGo = true
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

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == fmt.Sprintf("%v:%v", config.Domain, config.Port) {
			addr = fmt.Sprintf("%v:%v", config.Hosts[0], config.Port)
		}
		return dialer.DialContext(ctx, network, addr)
	}

	headers := http.Header{}
	headers.Add("User-Agent", "SitePingerDaemon/0.1")

	requestUrl, err := url.Parse(config.URI)
	if err != nil {
		log.Print(err)
	}

	request := http.Request{
		Method:           http.MethodGet,
		URL:              requestUrl,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           headers,
		Body:             nil,
		GetBody:          nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Host:             "",
		Form:             nil,
		PostForm:         nil,
		MultipartForm:    nil,
		Trailer:          nil,
		RemoteAddr:       "",
		RequestURI:       "",
		TLS:              nil,
		Cancel:           nil,
		Response:         nil,
	}

	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	resp, err := client.Do(&request)

	//resp, err := http.Get(config.URL)
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

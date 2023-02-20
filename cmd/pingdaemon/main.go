package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pinguinens/site-pinger/internal/config"
	"github.com/pinguinens/site-pinger/internal/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "", "Custom config path")
	flag.Parse()
}

func main() {
	appConf, err := config.New(configPath)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	appLogger, err := logger.New(appConf.GetLogFile())
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer logger.CloseLogFile()

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == fmt.Sprintf("%v:%v", appConf.Domain, appConf.Port) {
			addr = fmt.Sprintf("%v:%v", appConf.Hosts[0], appConf.Port)
		}
		return dialer.DialContext(ctx, network, addr)
	}

	headers := http.Header{}
	headers.Add("User-Agent", "SitePingerDaemon/0.1")

	requestUrl, err := url.Parse(appConf.URI)
	if err != nil {
		appLogger.Print(err)
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
	if err != nil {
		appLogger.Print(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		appLogger.Print(err)
	}

	appLogger.Info().Int("status_code", resp.StatusCode).Msg(string(body))
	appLogger.Debug().Msg("finish")
}

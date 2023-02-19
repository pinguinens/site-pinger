package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
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

	file, err := os.OpenFile(appConf.GetLogFile(), os.O_WRONLY|os.O_CREATE, 0755)
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
		if addr == fmt.Sprintf("%v:%v", appConf.Domain, appConf.Port) {
			addr = fmt.Sprintf("%v:%v", appConf.Hosts[0], appConf.Port)
		}
		return dialer.DialContext(ctx, network, addr)
	}

	headers := http.Header{}
	headers.Add("User-Agent", "SitePingerDaemon/0.1")

	requestUrl, err := url.Parse(appConf.URI)
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

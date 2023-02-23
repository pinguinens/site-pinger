package dialer

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var hosts HostTable

func New(timeout, keepAlive time.Duration) *net.Dialer {
	dialer := &net.Dialer{
		Timeout:   timeout * time.Second,
		KeepAlive: keepAlive * time.Second,
	}

	http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		parts := strings.Split(addr, ":")

		addr = fmt.Sprintf("%v:%v", hosts[parts[0]].Addr, parts[1])

		return dialer.DialContext(ctx, network, addr)
	}

	return dialer
}

func Hosts() {
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
}

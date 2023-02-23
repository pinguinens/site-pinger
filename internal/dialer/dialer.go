package dialer

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var resolveTable HostTable

func New(timeout, keepAlive time.Duration, hosts *HostTable) *net.Dialer {
	dialer := &net.Dialer{
		Timeout:   timeout * time.Second,
		KeepAlive: keepAlive * time.Second,
	}

	if hosts != nil {
		resolveTable = *hosts

		http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			parts := strings.Split(addr, ":")

			addr = fmt.Sprintf("%v:%v", resolveTable[parts[0]].Addr, parts[1])

			return dialer.DialContext(ctx, network, addr)
		}
	}

	return dialer
}

func Ping(uri string) (*http.Response, error) {
	client := http.Client{}

	headers := http.Header{}
	headers.Add("User-Agent", "SitePingerDaemon/0.1")

	requestUrl, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	request := http.Request{
		Method: http.MethodGet,
		URL:    requestUrl,
		Header: headers,
	}

	resp, err := client.Do(&request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

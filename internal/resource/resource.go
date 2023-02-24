package resource

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type Resource struct {
	Method string
	URI    string
	Host   Host
	dialer *net.Dialer
}

func New(method, uri string, host Host, dialer net.Dialer) Resource {
	return Resource{
		Method: method,
		URI:    uri,
		Host:   host,
		dialer: &dialer,
	}
}

func (r *Resource) Ping() (*http.Response, error) {
	dcFunc := func(ctx context.Context, network, addr string) (net.Conn, error) {
		parts := strings.Split(addr, ":")

		if addr == fmt.Sprintf("%v:%v", r.Host.Domain, parts[1]) {
			addr = fmt.Sprintf("%v:%v", r.Host.Addr, parts[1])
		}

		return r.dialer.DialContext(ctx, network, addr)
	}

	transport := &http.Transport{
		DialContext: dcFunc,
	}

	client := http.Client{
		Transport: transport,
	}

	headers := http.Header{}
	headers.Add("User-Agent", "SitePingerDaemon/0.1")

	requestUrl, err := url.Parse(r.URI)
	if err != nil {
		return nil, err
	}

	request := http.Request{
		Method: strings.ToUpper(r.Method),
		URL:    requestUrl,
		Header: headers,
	}

	resp, err := client.Do(&request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

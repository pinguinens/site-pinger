package connector

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type Connector struct {
	*http.Client
}

func New(domain, host string) *Connector {
	dialer := &net.Dialer{}

	transport := http.Transport{
		DialContext: makeDialContext(domain, host, dialer),
	}

	client := &http.Client{
		Transport: &transport,
	}

	return &Connector{client}
}

func makeDialContext(domain, host string, dialer *net.Dialer) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		parts := strings.Split(addr, ":")

		if addr == fmt.Sprintf("%v:%v", domain, parts[1]) {
			addr = fmt.Sprintf("%v:%v", host, parts[1])
		}

		return dialer.DialContext(ctx, network, addr)
	}
}

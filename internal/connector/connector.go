package connector

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/pinguinens/site-pinger/internal/site"
)

type Connector struct {
	*http.Client
}

func New(hosts site.HostTable) *Connector {
	dialer := &net.Dialer{}

	transport := http.Transport{
		DialContext: makeDialContext(dialer, hosts),
	}

	client := &http.Client{
		Transport: &transport,
	}

	return &Connector{client}
}

func makeDialContext(dialer *net.Dialer, hosts site.HostTable) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		reqUri := strings.Split(addr, ":")

		if ip, ok := hosts[reqUri[0]]; ok {
			addr = fmt.Sprintf("%v:%v", ip, reqUri[1])
		}

		return dialer.DialContext(ctx, network, addr)
	}
}

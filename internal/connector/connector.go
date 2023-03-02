package connector

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/pinguinens/site-pinger/internal/site"
)

func New(hosts site.HostTable) *http.Client {
	dialer := &net.Dialer{}

	transport := http.Transport{
		DialContext: makeDialContext(dialer, hosts),
	}

	client := &http.Client{
		Transport: &transport,
	}

	return client
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

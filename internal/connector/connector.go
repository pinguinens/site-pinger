package connector

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/pinguinens/site-pinger/internal/site"
)

const (
	clientName = "SitePinger"
)

type Connector struct {
	Client        *http.Client
	identificator string
}

func New(hosts site.HostTable, version string) *Connector {
	return &Connector{
		&http.Client{
			Transport: &http.Transport{
				DialContext: makeDialContext(&net.Dialer{}, hosts),
			},
		},
		fmt.Sprintf("%v/%v", clientName, version),
	}
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

func (c *Connector) GetIdentificator() string {
	return c.identificator
}

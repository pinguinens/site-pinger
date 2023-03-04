package resource

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-http-utils/headers"

	"github.com/pinguinens/site-pinger/internal/connector"
)

type Resource struct {
	Method    string
	URI       string
	Host      Host
	connector *connector.Connector
}

func New(method, uri string, host Host, client *connector.Connector) Resource {
	return Resource{
		Method:    method,
		URI:       uri,
		Host:      host,
		connector: client,
	}
}

func (r *Resource) Ping() (*http.Response, error) {
	reqHeaders := http.Header{}
	reqHeaders.Add(headers.UserAgent, r.connector.GetIdentificator())

	requestUrl, err := url.Parse(r.URI)
	if err != nil {
		return nil, err
	}

	request := http.Request{
		Method:     strings.ToUpper(r.Method),
		URL:        requestUrl,
		Header:     reqHeaders,
		RemoteAddr: r.Host.Addr,
	}

	resp, err := r.connector.Client.Do(&request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

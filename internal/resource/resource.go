package resource

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/pinguinens/site-pinger/internal/connector"
)

type Resource struct {
	Method    string
	URI       string
	Host      Host
	connector *connector.Connector
}

func New(method, uri string, host Host, connector *connector.Connector) Resource {
	return Resource{
		Method:    method,
		URI:       uri,
		Host:      host,
		connector: connector,
	}
}

func (r *Resource) Ping() (*http.Response, error) {
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

	resp, err := r.connector.Do(&request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

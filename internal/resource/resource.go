package resource

import (
	"net/http"
	"net/url"
	"strings"
)

type Resource struct {
	Method string
	URI    string
	Host   Host
	client *http.Client
}

func New(method, uri string, host Host, client *http.Client) Resource {
	return Resource{
		Method: method,
		URI:    uri,
		Host:   host,
		client: client,
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

	resp, err := r.client.Do(&request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

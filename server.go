package puppetdb

import (
	"net/http"
	"time"
)

/*
Server Representation of a PuppetDB server instance.

Use NewServer to create a new instance.
*/
type Server struct {
	BaseURL       string
	HTTPTransport http.RoundTripper
	HTTPTimeout   time.Duration
	Headers       map[string]string
}

// SetHTTPTimeout to set custom Timeout of http.Client
func (s *Server) SetHTTPTimeout(t time.Duration) {
	s.HTTPTimeout = t
}

// SetHeader the header
func (s *Server) SetHeader(key string, value string) {
	s.Headers[key] = value
}

/*
NewServer Create a new instance of a Server for usage later.

This is usually the main entry point of this SDK, where you would create
this initial object and use it to perform activities on the instance in
question.
*/
func NewServer(baseURL string) Server {
	return newServer(baseURL, nil)
}

/*
NewServerWithTransport Create a new instance of a Server for usage later.

Comparable to NewServer, but with an additional parameter to specify the http transport
(i.e. SSL options)
*/
func NewServerWithTransport(baseURL string, httpTransport http.RoundTripper) Server {
	return newServer(baseURL, httpTransport)
}

func newServer(baseURL string, httpTransport http.RoundTripper) Server {
	return Server{
		BaseURL:       baseURL,
		HTTPTransport: httpTransport,
		HTTPTimeout:   time.Second * 30,
		Headers:       make(map[string]string),
	}
}

package proxyerrors

import "errors"

var (
	ErrUpstreamClientError = errors.New("upstream server returned a client error (4xx)")
	ErrUpstreamServerError = errors.New("upstream server returned a server error (5xx)")
)

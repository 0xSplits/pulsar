package client

import "net/http"

type Interface interface {
	// Endpoint returns the configured base URL for the Neighbourhood API.
	Endpoint() string

	// Request performs an HTTP request to the Indexing Co Neighbourhood API.
	Request(*http.Request, any) error
}

package client

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xh3b4sd/tracer"
)

type Config struct {
	// Cli is the optional HTTP client being used to handle requests and
	// responses. Defaults to http.DefaultClient.
	Cli *http.Client

	// Url is the required API endpoint, e.g. "https://server.testing.splits.org".
	Url string
}

type Client struct {
	cli *http.Client
	url string
}

func New(c Config) *Client {
	if c.Cli == nil {
		c.Cli = http.DefaultClient
	}
	if c.Url == "" {
		tracer.Panic(tracer.Mask(fmt.Errorf("%T.Url must not be empty", c)))
	}

	return &Client{
		cli: c.Cli,
		url: strings.TrimSuffix(strings.TrimSpace(c.Url), "/"),
	}
}

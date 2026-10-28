package gothumbor

import (
	"net/http"
	"time"
)

const defaultTimeout = 15 * time.Second

// Client talks to a Thumbor server and signs image URLs with its secret.
type Client struct {
	httpClient *http.Client
	secret     string
	URL        string
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient replaces the default HTTP client (15s timeout).
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// NewClient returns a Client for the Thumbor server at url, signing paths with secret.
func NewClient(url, secret string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		secret:     secret,
		URL:        url,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

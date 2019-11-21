package gothumbor

import (
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	secret     string
	Url        string
}

func NewClient(url, secret string) *Client {
	client := Client{
		httpClient: &http.Client{Timeout: time.Second * 15},
		secret:     secret,
		Url:        url,
	}
	return &client
}

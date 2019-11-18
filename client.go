package gothumbor

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	secret     string
	url        string
}

func NewClient(host, port, secret string) *Client {
	url := fmt.Sprintf("http://%s:%s", host, port)

	client := Client{
		httpClient: &http.Client{Timeout: time.Second * 15},
		secret:     secret,
		url:        url,
	}
	return &client
}

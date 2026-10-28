package gothumbor

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

// ImageUpload posts body to the Thumbor /image endpoint using imageName as the Slug.
// The caller must close the returned response body.
func (c *Client) ImageUpload(ctx context.Context, imageName string, body io.Reader) (*http.Response, error) {
	endpoint, err := url.JoinPath(c.URL, "image")
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Slug", imageName)
	return c.httpClient.Do(req)
}

// ImageDelete deletes the image at imagePath from the Thumbor /image endpoint.
// The caller must close the returned response body.
func (c *Client) ImageDelete(ctx context.Context, imagePath string) (*http.Response, error) {
	endpoint, err := url.JoinPath(c.URL, "image", imagePath)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

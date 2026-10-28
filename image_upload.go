package gothumbor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ErrInvalidImagePath is returned when an image path could escape the /image endpoint.
var ErrInvalidImagePath = errors.New("gothumbor: invalid image path")

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

// ImageDelete deletes the image at imagePath, relative to the Thumbor /image endpoint.
// Paths with empty, "." or ".." segments are rejected with ErrInvalidImagePath.
// The caller must close the returned response body.
func (c *Client) ImageDelete(ctx context.Context, imagePath string) (*http.Response, error) {
	segments := strings.Split(strings.Trim(imagePath, "/"), "/")
	for _, s := range segments {
		if s == "" || s == "." || s == ".." {
			return nil, fmt.Errorf("%w: %q", ErrInvalidImagePath, imagePath)
		}
	}
	// #nosec G125 -- segments are validated above; escaping here would double-encode
	endpoint, err := url.JoinPath(c.URL, append([]string{"image"}, segments...)...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

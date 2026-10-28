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
// imagePath may be percent-encoded (as in Thumbor's Location header). Segments that
// decode to empty, ".", ".." or contain a slash or backslash are rejected with
// ErrInvalidImagePath. The caller must close the returned response body.
func (c *Client) ImageDelete(ctx context.Context, imagePath string) (*http.Response, error) {
	segments := strings.Split(strings.Trim(imagePath, "/"), "/")
	escaped := make([]string, 0, len(segments)+1)
	escaped = append(escaped, "image")
	for _, s := range segments {
		decoded, err := url.PathUnescape(s)
		if err != nil || decoded == "" || decoded == "." || decoded == ".." || strings.ContainsAny(decoded, `/\`) {
			return nil, fmt.Errorf("%w: %q", ErrInvalidImagePath, imagePath)
		}
		escaped = append(escaped, url.PathEscape(decoded))
	}
	// #nosec G125 -- every element is decoded, validated and re-escaped with url.PathEscape above
	endpoint, err := url.JoinPath(c.URL, escaped...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

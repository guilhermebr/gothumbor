package gothumbor_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/guilhermebr/gothumbor"
)

type capturedRequest struct {
	method string
	path   string
	slug   string
	body   string
}

func newTestServer(t *testing.T, status int) (*httptest.Server, *capturedRequest) {
	t.Helper()
	captured := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		captured.method = r.Method
		captured.path = r.URL.Path
		captured.slug = r.Header.Get("Slug")
		captured.body = string(body)
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	return srv, captured
}

func closeBody(t *testing.T, resp *http.Response) {
	t.Helper()
	if resp != nil {
		assert.NoError(t, resp.Body.Close())
	}
}

func TestNewClient(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		c := gothumbor.NewClient("http://thumbor.local", "secret")
		assert.Equal(t, "http://thumbor.local", c.URL)
	})

	t.Run("with http client", func(t *testing.T) {
		hc := &http.Client{Timeout: time.Second}
		srv, _ := newTestServer(t, http.StatusOK)
		c := gothumbor.NewClient(srv.URL, "secret", gothumbor.WithHTTPClient(hc))

		resp, err := c.ImageDelete(context.Background(), "abc")
		require.NoError(t, err)
		defer closeBody(t, resp)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestClientImageUpload(t *testing.T) {
	tests := []struct {
		name   string
		base   string
		status int
	}{
		{name: "created", base: "", status: http.StatusCreated},
		{name: "base url with trailing slash", base: "/", status: http.StatusCreated},
		{name: "server error is returned as response", base: "", status: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, captured := newTestServer(t, tt.status)
			c := gothumbor.NewClient(srv.URL+tt.base, "secret")

			resp, err := c.ImageUpload(context.Background(), "photo.jpg", strings.NewReader("image-bytes"))
			require.NoError(t, err)
			defer closeBody(t, resp)

			assert.Equal(t, tt.status, resp.StatusCode)
			assert.Equal(t, http.MethodPost, captured.method)
			assert.Equal(t, "/image", captured.path)
			assert.Equal(t, "photo.jpg", captured.slug)
			assert.Equal(t, "image-bytes", captured.body)
		})
	}
}

func TestClientImageDelete(t *testing.T) {
	tests := []struct {
		name      string
		imagePath string
		wantPath  string
	}{
		{name: "simple id", imagePath: "abc123", wantPath: "/image/abc123"},
		{name: "nested path", imagePath: "abc123/photo.jpg", wantPath: "/image/abc123/photo.jpg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, captured := newTestServer(t, http.StatusNoContent)
			c := gothumbor.NewClient(srv.URL, "secret")

			resp, err := c.ImageDelete(context.Background(), tt.imagePath)
			require.NoError(t, err)
			defer closeBody(t, resp)

			assert.Equal(t, http.StatusNoContent, resp.StatusCode)
			assert.Equal(t, http.MethodDelete, captured.method)
			assert.Equal(t, tt.wantPath, captured.path)
		})
	}
}

func TestClientContextCanceled(t *testing.T) {
	srv, _ := newTestServer(t, http.StatusOK)
	c := gothumbor.NewClient(srv.URL, "secret")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	resp, err := c.ImageUpload(ctx, "photo.jpg", strings.NewReader("x"))
	closeBody(t, resp)
	require.ErrorIs(t, err, context.Canceled)

	resp, err = c.ImageDelete(ctx, "abc")
	closeBody(t, resp)
	require.ErrorIs(t, err, context.Canceled)
}

func TestClientGetCryptedThumborPath(t *testing.T) {
	c := gothumbor.NewClient("http://thumbor.local", myKey)

	got, err := c.GetCryptedThumborPath(imageURL, gothumbor.ThumborOptions{Width: width, Height: height})
	require.NoError(t, err)
	assert.Equal(t, widthHeightEncryptedURL, got)
}

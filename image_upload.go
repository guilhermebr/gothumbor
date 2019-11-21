package gothumbor

import (
	"fmt"
	"io"
	"net/http"
)

func (c *Client) ImageUpload(imageName string, body io.Reader) (resp *http.Response, err error) {
	path := fmt.Sprintf("%s/image", c.Url)
	request, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Slug", imageName)
	return c.httpClient.Do(request)
}

func (c *Client) ImageDelete(imagePath string) (resp *http.Response, err error) {
	path := fmt.Sprintf("%s/image/%s", c.Url, imagePath)
	request, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(request)
}

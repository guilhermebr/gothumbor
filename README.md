gothumbor
=========

[![Go Reference](https://pkg.go.dev/badge/github.com/guilhermebr/gothumbor.svg)](https://pkg.go.dev/github.com/guilhermebr/gothumbor)
[![ci](https://github.com/guilhermebr/gothumbor/actions/workflows/ci.yml/badge.svg)](https://github.com/guilhermebr/gothumbor/actions/workflows/ci.yml)

gothumbor allows easy usage of [Thumbor] in Go. Requires Go 1.27+.


Installation
------------
```sh
go get github.com/guilhermebr/gothumbor
```


Signing URLs
------------
```go
myKey := "my-very-secret-key"
myImageURL := "my-domain.com/static/images/fancy-image.png"
thumborOptions := gothumbor.ThumborOptions{Width: 540, Height: 480}
newURL, err := gothumbor.GetCryptedThumborPath(myKey, myImageURL, thumborOptions)
```

With these variables the values of:

* newURL: AGp4diIF89Cm2ugmDGjhycikYjY=/540x480/my-domain.com/static/images/fancy-image.png
* err: nil


Client
------
`Client` signs paths with its secret and talks to Thumbor's `/image` upload endpoint.

```go
client := gothumbor.NewClient("https://thumbor.example.com", "my-very-secret-key",
	gothumbor.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}), // optional, default 15s
)

path, err := client.GetCryptedThumborPath(myImageURL, gothumbor.ThumborOptions{Width: 540, Height: 480})

resp, err := client.ImageUpload(ctx, "fancy-image.png", file)
if err != nil {
	return err
}
defer resp.Body.Close()
location := resp.Header.Get("Location") // e.g. /image/<id>/fancy-image.png

resp, err = client.ImageDelete(ctx, "<id>/fancy-image.png")
```

The caller must close the response body. Non-2xx responses aren't treated as errors, so check `resp.StatusCode`.


Development
-----------
```sh
make setup  # install golangci-lint, gosec, govulncheck
make check  # vet, lint, test (-race), gosec, govulncheck
```


License
-------

[MIT][mit] © [Globo.com][globocom]


Contributors
------------

Click [here][contributors] to see the list of contributors.

[Thumbor]:      https://github.com/thumbor/thumbor
[mit]:          https://github.com/guilhermebr/gothumbor/blob/main/LICENSE
[globocom]:     https://github.com/globocom
[contributors]: https://github.com/guilhermebr/gothumbor/graphs/contributors

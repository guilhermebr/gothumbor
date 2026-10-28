package gothumbor

import (
	"crypto/hmac"
	"crypto/sha1" // #nosec G505 -- required by the Thumbor URL signing protocol
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// ThumborOptions describes the image operations encoded in a Thumbor path.
type ThumborOptions struct {
	Trim    bool
	Width   int
	Height  int
	VAlign  string
	Smart   bool
	FitIn   bool
	Filters []string
	Flip    bool
	Flop    bool
	Meta    bool
	Left    int
	Top     int
	Right   int
	Bottom  int
}

// GetCryptedThumborPath returns the signed Thumbor path for imageURL using the client secret.
func (c *Client) GetCryptedThumborPath(imageURL string, options ThumborOptions) (string, error) {
	return GetCryptedThumborPath(c.secret, imageURL, options)
}

// GetCryptedThumborPath returns the Thumbor path for imageURL signed with key.
func GetCryptedThumborPath(key, imageURL string, options ThumborOptions) (string, error) {
	partial, err := GetThumborPath(imageURL, options)
	if err != nil {
		return "", err
	}
	hash := hmac.New(sha1.New, []byte(key))
	hash.Write([]byte(partial))
	signature := base64.URLEncoding.EncodeToString(hash.Sum(nil))
	return signature + "/" + partial, nil
}

// GetThumborPath returns the unsigned Thumbor path for imageURL.
func GetThumborPath(imageURL string, options ThumborOptions) (string, error) {
	return getURLParts(imageURL, options)
}

func getURLParts(imageURL string, options ThumborOptions) (string, error) {
	var parts []string

	partialObject, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}
	newImageURL := partialObject.String()
	if options.Meta {
		parts = append(parts, "meta")
	}

	if options.Trim {
		parts = append(parts, "trim")
	}

	if options.FitIn {
		parts = append(parts, "fit-in")
	}

	if options.Left != 0 || options.Top != 0 || options.Right != 0 || options.Bottom != 0 {
		parts = append(parts, fmt.Sprintf("%dx%d:%dx%d", options.Left, options.Top, options.Right, options.Bottom))
	}

	if options.Height != 0 || options.Width != 0 || options.Flip || options.Flop {
		flip := ""
		flop := ""

		if options.Flip {
			flip = "-"
		}
		if options.Flop {
			flop = "-"
		}
		parts = append(parts, fmt.Sprintf("%s%dx%s%d", flip, options.Width, flop, options.Height))
	}

	if options.VAlign != "" {
		parts = append(parts, options.VAlign)
	}
	if options.Smart {
		parts = append(parts, "smart")
	}

	if len(options.Filters) > 0 {
		parts = append(parts, "filters:"+strings.Join(options.Filters, ":"))
	}

	parts = append(parts, newImageURL)
	return strings.Join(parts, "/"), nil
}

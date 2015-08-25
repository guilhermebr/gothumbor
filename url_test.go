package gothumbor_test

import (
	"strings"
	"testing"

	"github.com/globocom/gothumbor"
)

const MYKEY = "my-security-key"
const IMAGEURL = "my.server.com/some/path/to/image.jpg"
const WIDTH = 300
const HEIGHT = 200
const ENCRYPTED_URL = "8ammJH8D-7tXy6kU3lTvoXlhu4o=/300x200/my.server.com/some/path/to/image.jpg"

func TestGetUrlPartialWithDefautWidthAndHeight(t *testing.T) {
	thumborOptions := gothumbor.ThumborOptions{}
	url, err := gothumbor.GetUrlParts(IMAGEURL, thumborOptions)
	if err != nil {
		t.Errorf("Got an error when tried to generate the thumbor url", err)
	}
	if url != IMAGEURL {
		t.Error("Got an unxpected partial url:", url)
	}
}

func TestGetUrlPartialWithSmart(t *testing.T) {
	thumborOptions := gothumbor.ThumborOptions{Smart: true}
	url, err := gothumbor.GetUrlParts(IMAGEURL, thumborOptions)
	if err != nil {
		t.Errorf("Got an error when tried to generate the thumbor url", err)
	}
	if url != strings.Join([]string{"smart", IMAGEURL}, "/") {
		t.Error("Got an unxpected partial url:", url)
	}
}

func TestGetUrlPartialOnlyWithWidthAndHeight(t *testing.T) {
	thumborOptions := gothumbor.ThumborOptions{Width: WIDTH, Height: HEIGHT}
	url, err := gothumbor.GetUrlParts(IMAGEURL, thumborOptions)
	if err != nil {
		t.Errorf("Got an error when tried to generate the thumbor url", err)
	}
	if url != strings.Join([]string{"300x200", IMAGEURL}, "/") {
		t.Error("Got an unxpected partial path:", url)
	}
}

func TestGetUrlUnderSpec1(t *testing.T) {
	//For spec 1: https://github.com/thumbor/thumbor/wiki/Libraries

	thumborOptions := gothumbor.ThumborOptions{Width: WIDTH, Height: HEIGHT}

	newUrl, err := gothumbor.GetThumborPath(MYKEY, IMAGEURL, thumborOptions)
	if err != nil {
		t.Errorf("Got an error when tried to generate the thumbor url", err)
	}
	if newUrl != ENCRYPTED_URL {
		t.Error("Got an unxpected thumbor path:", newUrl)
	}
}

func TestNegativeHeight(t *testing.T) {
	thumborOptions := gothumbor.ThumborOptions{
		Height: -1,
	}
	url, err := gothumbor.GetUrlParts(IMAGEURL, thumborOptions)
	if err == nil || url != "" {
		t.Errorf("Got an error when tried to generate the thumbor url")
	}

	if err.Error() != gothumbor.ErrorHeight.Error() {
		t.Errorf("Got an unxpected error height")
	}
}

func TestNegativeWidth(t *testing.T) {
	thumborOptions := gothumbor.ThumborOptions{
		Width: -1,
	}
	url, err := gothumbor.GetUrlParts(IMAGEURL, thumborOptions)
	if err == nil || url != "" {
		t.Errorf("Got an error when tried to generate the thumbor url")
	}

	if err.Error() != gothumbor.ErrorWidth.Error() {
		t.Errorf("Got an unxpected error width")
	}
}

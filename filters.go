package gothumbor

import (
	"fmt"
	"net/url"
	"strings"
)

// DistributedCollageFilter builds a distributed_collage filter from the given image URLs.
func DistributedCollageFilter(orientation, alignment string, urls []string) string {
	return fmt.Sprintf(
		"distributed_collage(%s,%s,%s)",
		orientation,
		alignment,
		strings.Join(urls, url.QueryEscape("|")),
	)
}

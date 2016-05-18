// +build go1.6

package strings

import (
	"strings"
)

func init() {
	Exports["lastIndexByte"] = strings.LastIndexByte
	Exports["compare"] = strings.Compare
}

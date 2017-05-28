// +build go1.6

package io

import (
	"io"
)

func init() {
	Exports["copyBuffer"] = io.CopyBuffer
	Exports["CopyBuffer"] = io.CopyBuffer
}

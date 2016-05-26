package sha1

import (
	"crypto/sha1"
	"encoding/hex"
)

// -----------------------------------------------------------------------------

// Sumstr is hex.EncodeToString(sha1.Sum(b)).
//
func Sumstr(b []byte) string {

	v := sha1.Sum(b)
	return hex.EncodeToString(v[:])
}

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":  "crypto/sha1",
	"new":    sha1.New,
	"sum":    sha1.Sum,
	"sumstr": Sumstr,

	"BlockSize": sha1.BlockSize,
	"Size":      sha1.Size,
}

// -----------------------------------------------------------------------------

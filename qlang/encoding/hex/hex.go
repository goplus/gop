package hex

import "encoding/hex"

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":          "encoding/hex",
	"encodedLen":     hex.EncodedLen,
	"encode":         hex.Encode,
	"decodedLen":     hex.DecodedLen,
	"decode":         hex.Decode,
	"encodeToString": hex.EncodeToString,
	"decodeString":   hex.DecodeString,
	"dump":           hex.Dump,
	"dumper":         hex.Dumper,

	"EncodedLen":     hex.EncodedLen,
	"Encode":         hex.Encode,
	"DecodedLen":     hex.DecodedLen,
	"Decode":         hex.Decode,
	"EncodeToString": hex.EncodeToString,
	"DecodeString":   hex.DecodeString,
	"Dump":           hex.Dump,
	"Dumper":         hex.Dumper,
}

// -----------------------------------------------------------------------------

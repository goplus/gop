package lib

//go:generate go run ../cmd/qexp archive/...
//go:generate go run ../cmd/qexp bufio bytes
//go:generate go run ../cmd/qexp compress/... context crypto/...
//go:generate go run ../cmd/qexp database/...
//go:generate go run ../cmd/qexp encoding/... expvar
//go:generate go run ../cmd/qexp flag fmt
//go:generate go run ../cmd/qexp hash/... html/...
//go:generate go run ../cmd/qexp image/... io/...
//go:generate go run ../cmd/qexp log/...
//go:generate go run ../cmd/qexp math/... mime/...
//go:generate go run ../cmd/qexp net/...
//go:generate go run ../cmd/qexp os/...
//go:generate go run ../cmd/qexp path/...
//go:generate go run ../cmd/qexp runtime reflect
//go:generate go run ../cmd/qexp sort strconv sync/...
//go:generate go run ../cmd/qexp time
//go:generate go run ../cmd/qexp unicode/...
//go:generate go run ../cmd/qexp github.com/goplus/gop/ast/gopq github.com/goplus/gop/ast/goptest

import (
	_ "github.com/goplus/gop/lib/builtin" // builtin
	_ "github.com/goplus/gop/lib/unsafe"  // unsafe

	_ "github.com/goplus/gop/lib/archive/tar"
	_ "github.com/goplus/gop/lib/archive/zip"
	_ "github.com/goplus/gop/lib/bufio"
	_ "github.com/goplus/gop/lib/bytes"
	_ "github.com/goplus/gop/lib/compress/bzip2"
	_ "github.com/goplus/gop/lib/compress/flate"
	_ "github.com/goplus/gop/lib/compress/gzip"
	_ "github.com/goplus/gop/lib/compress/lzw"
	_ "github.com/goplus/gop/lib/compress/zlib"
	_ "github.com/goplus/gop/lib/context"
	_ "github.com/goplus/gop/lib/crypto"
	_ "github.com/goplus/gop/lib/crypto/aes"
	_ "github.com/goplus/gop/lib/crypto/cipher"
	_ "github.com/goplus/gop/lib/crypto/des"
	_ "github.com/goplus/gop/lib/crypto/dsa"
	_ "github.com/goplus/gop/lib/crypto/ecdsa"
	_ "github.com/goplus/gop/lib/crypto/ed25519"
	_ "github.com/goplus/gop/lib/crypto/elliptic"
	_ "github.com/goplus/gop/lib/crypto/hmac"
	_ "github.com/goplus/gop/lib/crypto/md5"
	_ "github.com/goplus/gop/lib/crypto/rand"
	_ "github.com/goplus/gop/lib/crypto/rc4"
	_ "github.com/goplus/gop/lib/crypto/rsa"
	_ "github.com/goplus/gop/lib/crypto/sha1"
	_ "github.com/goplus/gop/lib/crypto/sha256"
	_ "github.com/goplus/gop/lib/crypto/sha512"
	_ "github.com/goplus/gop/lib/crypto/subtle"
	_ "github.com/goplus/gop/lib/crypto/tls"
	_ "github.com/goplus/gop/lib/crypto/x509"
	_ "github.com/goplus/gop/lib/database/sql"
	_ "github.com/goplus/gop/lib/database/sql/driver"
	_ "github.com/goplus/gop/lib/encoding"
	_ "github.com/goplus/gop/lib/encoding/ascii85"
	_ "github.com/goplus/gop/lib/encoding/asn1"
	_ "github.com/goplus/gop/lib/encoding/base32"
	_ "github.com/goplus/gop/lib/encoding/base64"
	_ "github.com/goplus/gop/lib/encoding/binary"
	_ "github.com/goplus/gop/lib/encoding/csv"
	_ "github.com/goplus/gop/lib/encoding/gob"
	_ "github.com/goplus/gop/lib/encoding/hex"
	_ "github.com/goplus/gop/lib/encoding/json"
	_ "github.com/goplus/gop/lib/encoding/pem"
	_ "github.com/goplus/gop/lib/encoding/xml"
	_ "github.com/goplus/gop/lib/errors"
	_ "github.com/goplus/gop/lib/expvar"
	_ "github.com/goplus/gop/lib/flag"
	_ "github.com/goplus/gop/lib/fmt"
	_ "github.com/goplus/gop/lib/hash"
	_ "github.com/goplus/gop/lib/hash/adler32"
	_ "github.com/goplus/gop/lib/hash/crc32"
	_ "github.com/goplus/gop/lib/hash/crc64"
	_ "github.com/goplus/gop/lib/hash/fnv"
	_ "github.com/goplus/gop/lib/hash/maphash"
	_ "github.com/goplus/gop/lib/html"
	_ "github.com/goplus/gop/lib/html/template"
	_ "github.com/goplus/gop/lib/image"
	_ "github.com/goplus/gop/lib/image/color"
	_ "github.com/goplus/gop/lib/image/color/palette"
	_ "github.com/goplus/gop/lib/image/draw"
	_ "github.com/goplus/gop/lib/image/gif"
	_ "github.com/goplus/gop/lib/image/jpeg"
	_ "github.com/goplus/gop/lib/image/png"
	_ "github.com/goplus/gop/lib/io"
	_ "github.com/goplus/gop/lib/io/ioutil"
	_ "github.com/goplus/gop/lib/log"
	_ "github.com/goplus/gop/lib/math"
	_ "github.com/goplus/gop/lib/math/big"
	_ "github.com/goplus/gop/lib/math/bits"
	_ "github.com/goplus/gop/lib/math/cmplx"
	_ "github.com/goplus/gop/lib/math/rand"
	_ "github.com/goplus/gop/lib/mime"
	_ "github.com/goplus/gop/lib/mime/multipart"
	_ "github.com/goplus/gop/lib/mime/quotedprintable"
	_ "github.com/goplus/gop/lib/net"
	_ "github.com/goplus/gop/lib/net/http"
	_ "github.com/goplus/gop/lib/net/http/cgi"
	_ "github.com/goplus/gop/lib/net/http/cookiejar"
	_ "github.com/goplus/gop/lib/net/http/fcgi"
	_ "github.com/goplus/gop/lib/net/http/httptest"
	_ "github.com/goplus/gop/lib/net/http/httptrace"
	_ "github.com/goplus/gop/lib/net/http/pprof"
	_ "github.com/goplus/gop/lib/net/mail"
	_ "github.com/goplus/gop/lib/net/rpc"
	_ "github.com/goplus/gop/lib/net/rpc/jsonrpc"
	_ "github.com/goplus/gop/lib/net/smtp"
	_ "github.com/goplus/gop/lib/net/textproto"
	_ "github.com/goplus/gop/lib/net/url"
	_ "github.com/goplus/gop/lib/os"
	_ "github.com/goplus/gop/lib/os/exec"
	_ "github.com/goplus/gop/lib/os/signal"
	_ "github.com/goplus/gop/lib/os/user"
	_ "github.com/goplus/gop/lib/path"
	_ "github.com/goplus/gop/lib/path/filepath"
	_ "github.com/goplus/gop/lib/reflect"
	_ "github.com/goplus/gop/lib/runtime"
	_ "github.com/goplus/gop/lib/sort"
	_ "github.com/goplus/gop/lib/strconv"
	_ "github.com/goplus/gop/lib/strings"
	_ "github.com/goplus/gop/lib/sync"
	_ "github.com/goplus/gop/lib/sync/atomic"
	_ "github.com/goplus/gop/lib/time"
	_ "github.com/goplus/gop/lib/unicode"
	_ "github.com/goplus/gop/lib/unicode/utf16"
	_ "github.com/goplus/gop/lib/unicode/utf8"

	_ "github.com/goplus/gop/lib/github.com/goplus/gop/ast/gopq"
	_ "github.com/goplus/gop/lib/github.com/goplus/gop/ast/goptest"
)

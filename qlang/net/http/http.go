package http

import (
	"net/http"

	qlang "qlang.io/spec"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":            "net/http",
	"request":          http.NewRequest,
	"readRequest":      http.ReadRequest,
	"readResponse":     http.ReadResponse,
	"parseHTTPVersion": http.ParseHTTPVersion,
	"parseTime":        http.ParseTime,

	"get":      http.Get,
	"post":     http.Post,
	"postForm": http.PostForm,
	"head":     http.Head,

	"handle":            http.Handle,
	"handleFunc":        http.HandleFunc,
	"serveMux":          http.NewServeMux,
	"serve":             http.Serve,
	"listenAndServe":    http.ListenAndServe,
	"listenAndServeTLS": http.ListenAndServeTLS,

	"error":           http.Error,
	"notFound":        http.NotFound,
	"notFoundHandler": http.NotFoundHandler,
	"redirect":        http.Redirect,
	"redirectHandler": http.RedirectHandler,
	"fileServer":      http.FileServer,
	"fileTransport":   http.NewFileTransport,
	"serveContent":    http.ServeContent,
	"serveFile":       http.ServeFile,
	"setCookie":       http.SetCookie,
	"stripPrefix":     http.StripPrefix,
	"timeoutHandler":  http.TimeoutHandler,

	"statusText":           http.StatusText,
	"canonicalHeaderKey":   http.CanonicalHeaderKey,
	"detectContentType":    http.DetectContentType,
	"maxBytesReader":       http.MaxBytesReader,
	"proxyFromEnvironment": http.ProxyFromEnvironment,
	"proxyURL":             http.ProxyURL,

	"DefaultTransport": http.DefaultTransport,
	"DefaultClient":    http.DefaultClient,
	"DefaultServeMux":  http.DefaultServeMux,

	"Client":   qlang.StructOf((*http.Client)(nil)),
	"Cookie":   qlang.StructOf((*http.Cookie)(nil)),
	"Header":   qlang.StructOf((*http.Header)(nil)),
	"Request":  qlang.StructOf((*http.Request)(nil)),
	"Response": qlang.StructOf((*http.Response)(nil)),
	"Server":   qlang.StructOf((*http.Server)(nil)),
}

// -----------------------------------------------------------------------------

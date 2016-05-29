package http

import (
	"net/http"
	"reflect"

	"qlang.io/qlang.spec.v1"
)

// -----------------------------------------------------------------------------

var (
	request = qlang.NewTypeEx(reflect.TypeOf((*http.Request)(nil)).Elem(), http.NewRequest)
)

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":            "net/http",
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

	"Client":   qlang.NewType(reflect.TypeOf((*http.Client)(nil)).Elem()),
	"Cookie":   qlang.NewType(reflect.TypeOf((*http.Cookie)(nil)).Elem()),
	"Header":   qlang.NewType(reflect.TypeOf((*http.Header)(nil)).Elem()),
	"Response": qlang.NewType(reflect.TypeOf((*http.Response)(nil)).Elem()),
	"Server":   qlang.NewType(reflect.TypeOf((*http.Server)(nil)).Elem()),

	"Request": request,
	"request": request,
}

// -----------------------------------------------------------------------------

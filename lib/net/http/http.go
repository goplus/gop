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

	"NewRequest":       http.NewRequest,
	"ReadRequest":      http.ReadRequest,
	"ReadResponse":     http.ReadResponse,
	"ParseHTTPVersion": http.ParseHTTPVersion,
	"ParseTime":        http.ParseTime,

	"get":      http.Get,
	"post":     http.Post,
	"postForm": http.PostForm,
	"head":     http.Head,

	"Get":      http.Get,
	"Post":     http.Post,
	"PostForm": http.PostForm,
	"Head":     http.Head,

	"handle":            http.Handle,
	"handleFunc":        http.HandleFunc,
	"serveMux":          http.NewServeMux,
	"serve":             http.Serve,
	"listenAndServe":    http.ListenAndServe,
	"listenAndServeTLS": http.ListenAndServeTLS,

	"Handle":            http.Handle,
	"HandleFunc":        http.HandleFunc,
	"NewServeMux":       http.NewServeMux,
	"Serve":             http.Serve,
	"ListenAndServe":    http.ListenAndServe,
	"ListenAndServeTLS": http.ListenAndServeTLS,

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

	"Error":            http.Error,
	"NotFound":         http.NotFound,
	"NotFoundHandler":  http.NotFoundHandler,
	"Redirect":         http.Redirect,
	"RedirectHandler":  http.RedirectHandler,
	"FileServer":       http.FileServer,
	"NewFileTransport": http.NewFileTransport,
	"ServeContent":     http.ServeContent,
	"ServeFile":        http.ServeFile,
	"SetCookie":        http.SetCookie,
	"StripPrefix":      http.StripPrefix,
	"TimeoutHandler":   http.TimeoutHandler,

	"statusText":           http.StatusText,
	"canonicalHeaderKey":   http.CanonicalHeaderKey,
	"detectContentType":    http.DetectContentType,
	"maxBytesReader":       http.MaxBytesReader,
	"proxyFromEnvironment": http.ProxyFromEnvironment,
	"proxyURL":             http.ProxyURL,

	"StatusText":           http.StatusText,
	"CanonicalHeaderKey":   http.CanonicalHeaderKey,
	"DetectContentType":    http.DetectContentType,
	"MaxBytesReader":       http.MaxBytesReader,
	"ProxyFromEnvironment": http.ProxyFromEnvironment,
	"ProxyURL":             http.ProxyURL,

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

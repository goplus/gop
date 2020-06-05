package http

import "net/http"

// -----------------------------------------------------------------------------

func newCookie() *http.Cookie {
	return new(http.Cookie)
}

func newHeader() http.Header {
	return make(http.Header)
}

func newClient() *http.Client {
	return new(http.Client)
}

func newServer() *http.Server {
	return new(http.Server)
}

func newResponse() *http.Response {
	return new(http.Response)
}

// -----------------------------------------------------------------------------

type caller interface {
	Call(...interface{}) interface{}
}

func handleFunc(pattern string, f caller) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		f.Call(w, req)
	})
}

// -----------------------------------------------------------------------------

type serveMux struct {
	*http.ServeMux
}

func (p serveMux) HandleFunc(pattern string, f caller) {
	p.ServeMux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		f.Call(w, req)
	})
}

func newServeMux() serveMux {
	return serveMux{ServeMux: http.NewServeMux()}
}

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"request":          http.NewRequest,
	"response":         newResponse,
	"header":           newHeader,
	"client":           newClient,
	"server":           newServer,
	"readRequest":      http.ReadRequest,
	"readResponse":     http.ReadResponse,
	"parseHTTPVersion": http.ParseHTTPVersion,
	"parseTime":        http.ParseTime,

	"get":      http.Get,
	"post":     http.Post,
	"postForm": http.PostForm,
	"head":     http.Head,

	"handle":            http.Handle,
	"handleFunc":        handleFunc,
	"serveMux":          newServeMux,
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
	"cookie":          newCookie,
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
}

// -----------------------------------------------------------------------------

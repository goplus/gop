// Package http provide Go+ "net/http" package, as "net/http" package in Go.
package http

import (
	bufio "bufio"
	context "context"
	io "io"
	net "net"
	http "net/http"
	url "net/url"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execCanonicalHeaderKey(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := http.CanonicalHeaderKey(args[0].(string))
	p.Ret(1, ret0)
}

func execmClientGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*http.Client).Get(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmClientDo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*http.Client).Do(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func toType0(v interface{}) io.Reader {
	if v == nil {
		return nil
	}
	return v.(io.Reader)
}

func execmClientPost(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := args[0].(*http.Client).Post(args[1].(string), args[2].(string), toType0(args[3]))
	p.Ret(4, ret0, ret1)
}

func execmClientPostForm(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(*http.Client).PostForm(args[1].(string), args[2].(url.Values))
	p.Ret(3, ret0, ret1)
}

func execmClientHead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*http.Client).Head(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmClientCloseIdleConnections(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*http.Client).CloseIdleConnections()
	p.PopN(1)
}

func execiCloseNotifierCloseNotify(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(http.CloseNotifier).CloseNotify()
	p.Ret(1, ret0)
}

func execmConnStateString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(http.ConnState).String()
	p.Ret(1, ret0)
}

func execmCookieString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Cookie).String()
	p.Ret(1, ret0)
}

func execiCookieJarCookies(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(http.CookieJar).Cookies(args[1].(*url.URL))
	p.Ret(2, ret0)
}

func execiCookieJarSetCookies(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(http.CookieJar).SetCookies(args[1].(*url.URL), args[2].([]*http.Cookie))
	p.PopN(3)
}

func execDetectContentType(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := http.DetectContentType(args[0].([]byte))
	p.Ret(1, ret0)
}

func execmDirOpen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(http.Dir).Open(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func toType1(v interface{}) http.ResponseWriter {
	if v == nil {
		return nil
	}
	return v.(http.ResponseWriter)
}

func execError(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	http.Error(toType1(args[0]), args[1].(string), args[2].(int))
	p.PopN(3)
}

func execiFileClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(http.File).Close()
	p.Ret(1, ret0)
}

func execiFileRead(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(http.File).Read(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiFileReaddir(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(http.File).Readdir(args[1].(int))
	p.Ret(2, ret0, ret1)
}

func execiFileSeek(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := args[0].(http.File).Seek(args[1].(int64), args[2].(int))
	p.Ret(3, ret0, ret1)
}

func execiFileStat(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(http.File).Stat()
	p.Ret(1, ret0, ret1)
}

func toType2(v interface{}) http.FileSystem {
	if v == nil {
		return nil
	}
	return v.(http.FileSystem)
}

func execFileServer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := http.FileServer(toType2(args[0]))
	p.Ret(1, ret0)
}

func execiFileSystemOpen(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(http.FileSystem).Open(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execiFlusherFlush(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(http.Flusher).Flush()
	p.PopN(1)
}

func execGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := http.Get(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func toType3(v interface{}) http.Handler {
	if v == nil {
		return nil
	}
	return v.(http.Handler)
}

func execHandle(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	http.Handle(args[0].(string), toType3(args[1]))
	p.PopN(2)
}

func execHandleFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	http.HandleFunc(args[0].(string), args[1].(func(http.ResponseWriter, *http.Request)))
	p.PopN(2)
}

func execiHandlerServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(http.Handler).ServeHTTP(toType1(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmHandlerFuncServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(http.HandlerFunc).ServeHTTP(toType1(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execHead(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := http.Head(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmHeaderAdd(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(http.Header).Add(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmHeaderSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(http.Header).Set(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmHeaderGet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(http.Header).Get(args[1].(string))
	p.Ret(2, ret0)
}

func execmHeaderValues(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(http.Header).Values(args[1].(string))
	p.Ret(2, ret0)
}

func execmHeaderDel(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(http.Header).Del(args[1].(string))
	p.PopN(2)
}

func toType4(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmHeaderWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(http.Header).Write(toType4(args[1]))
	p.Ret(2, ret0)
}

func execmHeaderClone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(http.Header).Clone()
	p.Ret(1, ret0)
}

func execmHeaderWriteSubset(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(http.Header).WriteSubset(toType4(args[1]), args[2].(map[string]bool))
	p.Ret(3, ret0)
}

func execiHijackerHijack(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(http.Hijacker).Hijack()
	p.Ret(1, ret0, ret1, ret2)
}

func execListenAndServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := http.ListenAndServe(args[0].(string), toType3(args[1]))
	p.Ret(2, ret0)
}

func execListenAndServeTLS(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := http.ListenAndServeTLS(args[0].(string), args[1].(string), args[2].(string), toType3(args[3]))
	p.Ret(4, ret0)
}

func toType5(v interface{}) io.ReadCloser {
	if v == nil {
		return nil
	}
	return v.(io.ReadCloser)
}

func execMaxBytesReader(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := http.MaxBytesReader(toType1(args[0]), toType5(args[1]), args[2].(int64))
	p.Ret(3, ret0)
}

func execNewFileTransport(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := http.NewFileTransport(toType2(args[0]))
	p.Ret(1, ret0)
}

func execNewRequest(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := http.NewRequest(args[0].(string), args[1].(string), toType0(args[2]))
	p.Ret(3, ret0, ret1)
}

func toType6(v interface{}) context.Context {
	if v == nil {
		return nil
	}
	return v.(context.Context)
}

func execNewRequestWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0, ret1 := http.NewRequestWithContext(toType6(args[0]), args[1].(string), args[2].(string), toType0(args[3]))
	p.Ret(4, ret0, ret1)
}

func execNewServeMux(_ int, p *gop.Context) {
	ret0 := http.NewServeMux()
	p.Ret(0, ret0)
}

func execNotFound(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	http.NotFound(toType1(args[0]), args[1].(*http.Request))
	p.PopN(2)
}

func execNotFoundHandler(_ int, p *gop.Context) {
	ret0 := http.NotFoundHandler()
	p.Ret(0, ret0)
}

func execParseHTTPVersion(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := http.ParseHTTPVersion(args[0].(string))
	p.Ret(1, ret0, ret1, ret2)
}

func execParseTime(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := http.ParseTime(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execPost(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := http.Post(args[0].(string), args[1].(string), toType0(args[2]))
	p.Ret(3, ret0, ret1)
}

func execPostForm(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := http.PostForm(args[0].(string), args[1].(url.Values))
	p.Ret(2, ret0, ret1)
}

func execmProtocolErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.ProtocolError).Error()
	p.Ret(1, ret0)
}

func execProxyFromEnvironment(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := http.ProxyFromEnvironment(args[0].(*http.Request))
	p.Ret(1, ret0, ret1)
}

func execProxyURL(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := http.ProxyURL(args[0].(*url.URL))
	p.Ret(1, ret0)
}

func execiPusherPush(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(http.Pusher).Push(args[1].(string), args[2].(*http.PushOptions))
	p.Ret(3, ret0)
}

func execReadRequest(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := http.ReadRequest(args[0].(*bufio.Reader))
	p.Ret(1, ret0, ret1)
}

func execReadResponse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := http.ReadResponse(args[0].(*bufio.Reader), args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func execRedirect(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	http.Redirect(toType1(args[0]), args[1].(*http.Request), args[2].(string), args[3].(int))
	p.PopN(4)
}

func execRedirectHandler(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := http.RedirectHandler(args[0].(string), args[1].(int))
	p.Ret(2, ret0)
}

func execmRequestContext(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Request).Context()
	p.Ret(1, ret0)
}

func execmRequestWithContext(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).WithContext(toType6(args[1]))
	p.Ret(2, ret0)
}

func execmRequestClone(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).Clone(toType6(args[1]))
	p.Ret(2, ret0)
}

func execmRequestProtoAtLeast(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*http.Request).ProtoAtLeast(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmRequestUserAgent(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Request).UserAgent()
	p.Ret(1, ret0)
}

func execmRequestCookies(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Request).Cookies()
	p.Ret(1, ret0)
}

func execmRequestCookie(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*http.Request).Cookie(args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execmRequestAddCookie(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*http.Request).AddCookie(args[1].(*http.Cookie))
	p.PopN(2)
}

func execmRequestReferer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Request).Referer()
	p.Ret(1, ret0)
}

func execmRequestMultipartReader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*http.Request).MultipartReader()
	p.Ret(1, ret0, ret1)
}

func execmRequestWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).Write(toType4(args[1]))
	p.Ret(2, ret0)
}

func execmRequestWriteProxy(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).WriteProxy(toType4(args[1]))
	p.Ret(2, ret0)
}

func execmRequestBasicAuth(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(*http.Request).BasicAuth()
	p.Ret(1, ret0, ret1, ret2)
}

func execmRequestSetBasicAuth(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*http.Request).SetBasicAuth(args[1].(string), args[2].(string))
	p.PopN(3)
}

func execmRequestParseForm(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Request).ParseForm()
	p.Ret(1, ret0)
}

func execmRequestParseMultipartForm(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).ParseMultipartForm(args[1].(int64))
	p.Ret(2, ret0)
}

func execmRequestFormValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).FormValue(args[1].(string))
	p.Ret(2, ret0)
}

func execmRequestPostFormValue(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Request).PostFormValue(args[1].(string))
	p.Ret(2, ret0)
}

func execmRequestFormFile(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1, ret2 := args[0].(*http.Request).FormFile(args[1].(string))
	p.Ret(2, ret0, ret1, ret2)
}

func execmResponseCookies(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Response).Cookies()
	p.Ret(1, ret0)
}

func execmResponseLocation(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(*http.Response).Location()
	p.Ret(1, ret0, ret1)
}

func execmResponseProtoAtLeast(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*http.Response).ProtoAtLeast(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmResponseWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Response).Write(toType4(args[1]))
	p.Ret(2, ret0)
}

func execiResponseWriterHeader(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(http.ResponseWriter).Header()
	p.Ret(1, ret0)
}

func execiResponseWriterWrite(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(http.ResponseWriter).Write(args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execiResponseWriterWriteHeader(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(http.ResponseWriter).WriteHeader(args[1].(int))
	p.PopN(2)
}

func execiRoundTripperRoundTrip(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(http.RoundTripper).RoundTrip(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func toType7(v interface{}) net.Listener {
	if v == nil {
		return nil
	}
	return v.(net.Listener)
}

func execServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := http.Serve(toType7(args[0]), toType3(args[1]))
	p.Ret(2, ret0)
}

func toType8(v interface{}) io.ReadSeeker {
	if v == nil {
		return nil
	}
	return v.(io.ReadSeeker)
}

func execServeContent(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	http.ServeContent(toType1(args[0]), args[1].(*http.Request), args[2].(string), args[3].(time.Time), toType8(args[4]))
	p.PopN(5)
}

func execServeFile(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	http.ServeFile(toType1(args[0]), args[1].(*http.Request), args[2].(string))
	p.PopN(3)
}

func execmServeMuxHandler(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*http.ServeMux).Handler(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func execmServeMuxServeHTTP(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*http.ServeMux).ServeHTTP(toType1(args[1]), args[2].(*http.Request))
	p.PopN(3)
}

func execmServeMuxHandle(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*http.ServeMux).Handle(args[1].(string), toType3(args[2]))
	p.PopN(3)
}

func execmServeMuxHandleFunc(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*http.ServeMux).HandleFunc(args[1].(string), args[2].(func(http.ResponseWriter, *http.Request)))
	p.PopN(3)
}

func execServeTLS(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := http.ServeTLS(toType7(args[0]), toType3(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

func execmServerClose(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Server).Close()
	p.Ret(1, ret0)
}

func execmServerShutdown(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Server).Shutdown(toType6(args[1]))
	p.Ret(2, ret0)
}

func execmServerRegisterOnShutdown(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*http.Server).RegisterOnShutdown(args[1].(func()))
	p.PopN(2)
}

func execmServerListenAndServe(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Server).ListenAndServe()
	p.Ret(1, ret0)
}

func execmServerServe(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*http.Server).Serve(toType7(args[1]))
	p.Ret(2, ret0)
}

func execmServerServeTLS(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*http.Server).ServeTLS(toType7(args[1]), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

func execmServerSetKeepAlivesEnabled(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*http.Server).SetKeepAlivesEnabled(args[1].(bool))
	p.PopN(2)
}

func execmServerListenAndServeTLS(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*http.Server).ListenAndServeTLS(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execSetCookie(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	http.SetCookie(toType1(args[0]), args[1].(*http.Cookie))
	p.PopN(2)
}

func execStatusText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := http.StatusText(args[0].(int))
	p.Ret(1, ret0)
}

func execStripPrefix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := http.StripPrefix(args[0].(string), toType3(args[1]))
	p.Ret(2, ret0)
}

func execTimeoutHandler(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := http.TimeoutHandler(toType3(args[0]), args[1].(time.Duration), args[2].(string))
	p.Ret(3, ret0)
}

func execmTransportRoundTrip(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := args[0].(*http.Transport).RoundTrip(args[1].(*http.Request))
	p.Ret(2, ret0, ret1)
}

func execmTransportClone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*http.Transport).Clone()
	p.Ret(1, ret0)
}

func toType9(v interface{}) http.RoundTripper {
	if v == nil {
		return nil
	}
	return v.(http.RoundTripper)
}

func execmTransportRegisterProtocol(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*http.Transport).RegisterProtocol(args[1].(string), toType9(args[2]))
	p.PopN(3)
}

func execmTransportCloseIdleConnections(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*http.Transport).CloseIdleConnections()
	p.PopN(1)
}

func execmTransportCancelRequest(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*http.Transport).CancelRequest(args[1].(*http.Request))
	p.PopN(2)
}

// I is a Go package instance.
var I = gop.NewGoPackage("net/http")

func init() {
	I.RegisterFuncs(
		I.Func("CanonicalHeaderKey", http.CanonicalHeaderKey, execCanonicalHeaderKey),
		I.Func("(*Client).Get", (*http.Client).Get, execmClientGet),
		I.Func("(*Client).Do", (*http.Client).Do, execmClientDo),
		I.Func("(*Client).Post", (*http.Client).Post, execmClientPost),
		I.Func("(*Client).PostForm", (*http.Client).PostForm, execmClientPostForm),
		I.Func("(*Client).Head", (*http.Client).Head, execmClientHead),
		I.Func("(*Client).CloseIdleConnections", (*http.Client).CloseIdleConnections, execmClientCloseIdleConnections),
		I.Func("(CloseNotifier).CloseNotify", (http.CloseNotifier).CloseNotify, execiCloseNotifierCloseNotify),
		I.Func("(ConnState).String", (http.ConnState).String, execmConnStateString),
		I.Func("(*Cookie).String", (*http.Cookie).String, execmCookieString),
		I.Func("(CookieJar).Cookies", (http.CookieJar).Cookies, execiCookieJarCookies),
		I.Func("(CookieJar).SetCookies", (http.CookieJar).SetCookies, execiCookieJarSetCookies),
		I.Func("DetectContentType", http.DetectContentType, execDetectContentType),
		I.Func("(Dir).Open", (http.Dir).Open, execmDirOpen),
		I.Func("Error", http.Error, execError),
		I.Func("(File).Close", (http.File).Close, execiFileClose),
		I.Func("(File).Read", (http.File).Read, execiFileRead),
		I.Func("(File).Readdir", (http.File).Readdir, execiFileReaddir),
		I.Func("(File).Seek", (http.File).Seek, execiFileSeek),
		I.Func("(File).Stat", (http.File).Stat, execiFileStat),
		I.Func("FileServer", http.FileServer, execFileServer),
		I.Func("(FileSystem).Open", (http.FileSystem).Open, execiFileSystemOpen),
		I.Func("(Flusher).Flush", (http.Flusher).Flush, execiFlusherFlush),
		I.Func("Get", http.Get, execGet),
		I.Func("Handle", http.Handle, execHandle),
		I.Func("HandleFunc", http.HandleFunc, execHandleFunc),
		I.Func("(Handler).ServeHTTP", (http.Handler).ServeHTTP, execiHandlerServeHTTP),
		I.Func("(HandlerFunc).ServeHTTP", (http.HandlerFunc).ServeHTTP, execmHandlerFuncServeHTTP),
		I.Func("Head", http.Head, execHead),
		I.Func("(Header).Add", (http.Header).Add, execmHeaderAdd),
		I.Func("(Header).Set", (http.Header).Set, execmHeaderSet),
		I.Func("(Header).Get", (http.Header).Get, execmHeaderGet),
		I.Func("(Header).Values", (http.Header).Values, execmHeaderValues),
		I.Func("(Header).Del", (http.Header).Del, execmHeaderDel),
		I.Func("(Header).Write", (http.Header).Write, execmHeaderWrite),
		I.Func("(Header).Clone", (http.Header).Clone, execmHeaderClone),
		I.Func("(Header).WriteSubset", (http.Header).WriteSubset, execmHeaderWriteSubset),
		I.Func("(Hijacker).Hijack", (http.Hijacker).Hijack, execiHijackerHijack),
		I.Func("ListenAndServe", http.ListenAndServe, execListenAndServe),
		I.Func("ListenAndServeTLS", http.ListenAndServeTLS, execListenAndServeTLS),
		I.Func("MaxBytesReader", http.MaxBytesReader, execMaxBytesReader),
		I.Func("NewFileTransport", http.NewFileTransport, execNewFileTransport),
		I.Func("NewRequest", http.NewRequest, execNewRequest),
		I.Func("NewRequestWithContext", http.NewRequestWithContext, execNewRequestWithContext),
		I.Func("NewServeMux", http.NewServeMux, execNewServeMux),
		I.Func("NotFound", http.NotFound, execNotFound),
		I.Func("NotFoundHandler", http.NotFoundHandler, execNotFoundHandler),
		I.Func("ParseHTTPVersion", http.ParseHTTPVersion, execParseHTTPVersion),
		I.Func("ParseTime", http.ParseTime, execParseTime),
		I.Func("Post", http.Post, execPost),
		I.Func("PostForm", http.PostForm, execPostForm),
		I.Func("(*ProtocolError).Error", (*http.ProtocolError).Error, execmProtocolErrorError),
		I.Func("ProxyFromEnvironment", http.ProxyFromEnvironment, execProxyFromEnvironment),
		I.Func("ProxyURL", http.ProxyURL, execProxyURL),
		I.Func("(Pusher).Push", (http.Pusher).Push, execiPusherPush),
		I.Func("ReadRequest", http.ReadRequest, execReadRequest),
		I.Func("ReadResponse", http.ReadResponse, execReadResponse),
		I.Func("Redirect", http.Redirect, execRedirect),
		I.Func("RedirectHandler", http.RedirectHandler, execRedirectHandler),
		I.Func("(*Request).Context", (*http.Request).Context, execmRequestContext),
		I.Func("(*Request).WithContext", (*http.Request).WithContext, execmRequestWithContext),
		I.Func("(*Request).Clone", (*http.Request).Clone, execmRequestClone),
		I.Func("(*Request).ProtoAtLeast", (*http.Request).ProtoAtLeast, execmRequestProtoAtLeast),
		I.Func("(*Request).UserAgent", (*http.Request).UserAgent, execmRequestUserAgent),
		I.Func("(*Request).Cookies", (*http.Request).Cookies, execmRequestCookies),
		I.Func("(*Request).Cookie", (*http.Request).Cookie, execmRequestCookie),
		I.Func("(*Request).AddCookie", (*http.Request).AddCookie, execmRequestAddCookie),
		I.Func("(*Request).Referer", (*http.Request).Referer, execmRequestReferer),
		I.Func("(*Request).MultipartReader", (*http.Request).MultipartReader, execmRequestMultipartReader),
		I.Func("(*Request).Write", (*http.Request).Write, execmRequestWrite),
		I.Func("(*Request).WriteProxy", (*http.Request).WriteProxy, execmRequestWriteProxy),
		I.Func("(*Request).BasicAuth", (*http.Request).BasicAuth, execmRequestBasicAuth),
		I.Func("(*Request).SetBasicAuth", (*http.Request).SetBasicAuth, execmRequestSetBasicAuth),
		I.Func("(*Request).ParseForm", (*http.Request).ParseForm, execmRequestParseForm),
		I.Func("(*Request).ParseMultipartForm", (*http.Request).ParseMultipartForm, execmRequestParseMultipartForm),
		I.Func("(*Request).FormValue", (*http.Request).FormValue, execmRequestFormValue),
		I.Func("(*Request).PostFormValue", (*http.Request).PostFormValue, execmRequestPostFormValue),
		I.Func("(*Request).FormFile", (*http.Request).FormFile, execmRequestFormFile),
		I.Func("(*Response).Cookies", (*http.Response).Cookies, execmResponseCookies),
		I.Func("(*Response).Location", (*http.Response).Location, execmResponseLocation),
		I.Func("(*Response).ProtoAtLeast", (*http.Response).ProtoAtLeast, execmResponseProtoAtLeast),
		I.Func("(*Response).Write", (*http.Response).Write, execmResponseWrite),
		I.Func("(ResponseWriter).Header", (http.ResponseWriter).Header, execiResponseWriterHeader),
		I.Func("(ResponseWriter).Write", (http.ResponseWriter).Write, execiResponseWriterWrite),
		I.Func("(ResponseWriter).WriteHeader", (http.ResponseWriter).WriteHeader, execiResponseWriterWriteHeader),
		I.Func("(RoundTripper).RoundTrip", (http.RoundTripper).RoundTrip, execiRoundTripperRoundTrip),
		I.Func("Serve", http.Serve, execServe),
		I.Func("ServeContent", http.ServeContent, execServeContent),
		I.Func("ServeFile", http.ServeFile, execServeFile),
		I.Func("(*ServeMux).Handler", (*http.ServeMux).Handler, execmServeMuxHandler),
		I.Func("(*ServeMux).ServeHTTP", (*http.ServeMux).ServeHTTP, execmServeMuxServeHTTP),
		I.Func("(*ServeMux).Handle", (*http.ServeMux).Handle, execmServeMuxHandle),
		I.Func("(*ServeMux).HandleFunc", (*http.ServeMux).HandleFunc, execmServeMuxHandleFunc),
		I.Func("ServeTLS", http.ServeTLS, execServeTLS),
		I.Func("(*Server).Close", (*http.Server).Close, execmServerClose),
		I.Func("(*Server).Shutdown", (*http.Server).Shutdown, execmServerShutdown),
		I.Func("(*Server).RegisterOnShutdown", (*http.Server).RegisterOnShutdown, execmServerRegisterOnShutdown),
		I.Func("(*Server).ListenAndServe", (*http.Server).ListenAndServe, execmServerListenAndServe),
		I.Func("(*Server).Serve", (*http.Server).Serve, execmServerServe),
		I.Func("(*Server).ServeTLS", (*http.Server).ServeTLS, execmServerServeTLS),
		I.Func("(*Server).SetKeepAlivesEnabled", (*http.Server).SetKeepAlivesEnabled, execmServerSetKeepAlivesEnabled),
		I.Func("(*Server).ListenAndServeTLS", (*http.Server).ListenAndServeTLS, execmServerListenAndServeTLS),
		I.Func("SetCookie", http.SetCookie, execSetCookie),
		I.Func("StatusText", http.StatusText, execStatusText),
		I.Func("StripPrefix", http.StripPrefix, execStripPrefix),
		I.Func("TimeoutHandler", http.TimeoutHandler, execTimeoutHandler),
		I.Func("(*Transport).RoundTrip", (*http.Transport).RoundTrip, execmTransportRoundTrip),
		I.Func("(*Transport).Clone", (*http.Transport).Clone, execmTransportClone),
		I.Func("(*Transport).RegisterProtocol", (*http.Transport).RegisterProtocol, execmTransportRegisterProtocol),
		I.Func("(*Transport).CloseIdleConnections", (*http.Transport).CloseIdleConnections, execmTransportCloseIdleConnections),
		I.Func("(*Transport).CancelRequest", (*http.Transport).CancelRequest, execmTransportCancelRequest),
	)
	I.RegisterVars(
		I.Var("DefaultClient", &http.DefaultClient),
		I.Var("DefaultServeMux", &http.DefaultServeMux),
		I.Var("DefaultTransport", &http.DefaultTransport),
		I.Var("ErrAbortHandler", &http.ErrAbortHandler),
		I.Var("ErrBodyNotAllowed", &http.ErrBodyNotAllowed),
		I.Var("ErrBodyReadAfterClose", &http.ErrBodyReadAfterClose),
		I.Var("ErrContentLength", &http.ErrContentLength),
		I.Var("ErrHandlerTimeout", &http.ErrHandlerTimeout),
		I.Var("ErrHeaderTooLong", &http.ErrHeaderTooLong),
		I.Var("ErrHijacked", &http.ErrHijacked),
		I.Var("ErrLineTooLong", &http.ErrLineTooLong),
		I.Var("ErrMissingBoundary", &http.ErrMissingBoundary),
		I.Var("ErrMissingContentLength", &http.ErrMissingContentLength),
		I.Var("ErrMissingFile", &http.ErrMissingFile),
		I.Var("ErrNoCookie", &http.ErrNoCookie),
		I.Var("ErrNoLocation", &http.ErrNoLocation),
		I.Var("ErrNotMultipart", &http.ErrNotMultipart),
		I.Var("ErrNotSupported", &http.ErrNotSupported),
		I.Var("ErrServerClosed", &http.ErrServerClosed),
		I.Var("ErrShortBody", &http.ErrShortBody),
		I.Var("ErrSkipAltProtocol", &http.ErrSkipAltProtocol),
		I.Var("ErrUnexpectedTrailer", &http.ErrUnexpectedTrailer),
		I.Var("ErrUseLastResponse", &http.ErrUseLastResponse),
		I.Var("ErrWriteAfterFlush", &http.ErrWriteAfterFlush),
		I.Var("LocalAddrContextKey", &http.LocalAddrContextKey),
		I.Var("NoBody", &http.NoBody),
		I.Var("ServerContextKey", &http.ServerContextKey),
	)
	I.RegisterTypes(
		I.Type("Client", reflect.TypeOf((*http.Client)(nil)).Elem()),
		I.Type("CloseNotifier", reflect.TypeOf((*http.CloseNotifier)(nil)).Elem()),
		I.Type("ConnState", reflect.TypeOf((*http.ConnState)(nil)).Elem()),
		I.Type("Cookie", reflect.TypeOf((*http.Cookie)(nil)).Elem()),
		I.Type("CookieJar", reflect.TypeOf((*http.CookieJar)(nil)).Elem()),
		I.Type("Dir", reflect.TypeOf((*http.Dir)(nil)).Elem()),
		I.Type("File", reflect.TypeOf((*http.File)(nil)).Elem()),
		I.Type("FileSystem", reflect.TypeOf((*http.FileSystem)(nil)).Elem()),
		I.Type("Flusher", reflect.TypeOf((*http.Flusher)(nil)).Elem()),
		I.Type("Handler", reflect.TypeOf((*http.Handler)(nil)).Elem()),
		I.Type("HandlerFunc", reflect.TypeOf((*http.HandlerFunc)(nil)).Elem()),
		I.Type("Header", reflect.TypeOf((*http.Header)(nil)).Elem()),
		I.Type("Hijacker", reflect.TypeOf((*http.Hijacker)(nil)).Elem()),
		I.Type("ProtocolError", reflect.TypeOf((*http.ProtocolError)(nil)).Elem()),
		I.Type("PushOptions", reflect.TypeOf((*http.PushOptions)(nil)).Elem()),
		I.Type("Pusher", reflect.TypeOf((*http.Pusher)(nil)).Elem()),
		I.Type("Request", reflect.TypeOf((*http.Request)(nil)).Elem()),
		I.Type("Response", reflect.TypeOf((*http.Response)(nil)).Elem()),
		I.Type("ResponseWriter", reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()),
		I.Type("RoundTripper", reflect.TypeOf((*http.RoundTripper)(nil)).Elem()),
		I.Type("SameSite", reflect.TypeOf((*http.SameSite)(nil)).Elem()),
		I.Type("ServeMux", reflect.TypeOf((*http.ServeMux)(nil)).Elem()),
		I.Type("Server", reflect.TypeOf((*http.Server)(nil)).Elem()),
		I.Type("Transport", reflect.TypeOf((*http.Transport)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("DefaultMaxHeaderBytes", qspec.ConstUnboundInt, http.DefaultMaxHeaderBytes),
		I.Const("DefaultMaxIdleConnsPerHost", qspec.ConstUnboundInt, http.DefaultMaxIdleConnsPerHost),
		I.Const("MethodConnect", qspec.ConstBoundString, http.MethodConnect),
		I.Const("MethodDelete", qspec.ConstBoundString, http.MethodDelete),
		I.Const("MethodGet", qspec.ConstBoundString, http.MethodGet),
		I.Const("MethodHead", qspec.ConstBoundString, http.MethodHead),
		I.Const("MethodOptions", qspec.ConstBoundString, http.MethodOptions),
		I.Const("MethodPatch", qspec.ConstBoundString, http.MethodPatch),
		I.Const("MethodPost", qspec.ConstBoundString, http.MethodPost),
		I.Const("MethodPut", qspec.ConstBoundString, http.MethodPut),
		I.Const("MethodTrace", qspec.ConstBoundString, http.MethodTrace),
		I.Const("SameSiteDefaultMode", qspec.Int, http.SameSiteDefaultMode),
		I.Const("SameSiteLaxMode", qspec.Int, http.SameSiteLaxMode),
		I.Const("SameSiteNoneMode", qspec.Int, http.SameSiteNoneMode),
		I.Const("SameSiteStrictMode", qspec.Int, http.SameSiteStrictMode),
		I.Const("StateActive", qspec.Int, http.StateActive),
		I.Const("StateClosed", qspec.Int, http.StateClosed),
		I.Const("StateHijacked", qspec.Int, http.StateHijacked),
		I.Const("StateIdle", qspec.Int, http.StateIdle),
		I.Const("StateNew", qspec.Int, http.StateNew),
		I.Const("StatusAccepted", qspec.ConstUnboundInt, http.StatusAccepted),
		I.Const("StatusAlreadyReported", qspec.ConstUnboundInt, http.StatusAlreadyReported),
		I.Const("StatusBadGateway", qspec.ConstUnboundInt, http.StatusBadGateway),
		I.Const("StatusBadRequest", qspec.ConstUnboundInt, http.StatusBadRequest),
		I.Const("StatusConflict", qspec.ConstUnboundInt, http.StatusConflict),
		I.Const("StatusContinue", qspec.ConstUnboundInt, http.StatusContinue),
		I.Const("StatusCreated", qspec.ConstUnboundInt, http.StatusCreated),
		I.Const("StatusEarlyHints", qspec.ConstUnboundInt, http.StatusEarlyHints),
		I.Const("StatusExpectationFailed", qspec.ConstUnboundInt, http.StatusExpectationFailed),
		I.Const("StatusFailedDependency", qspec.ConstUnboundInt, http.StatusFailedDependency),
		I.Const("StatusForbidden", qspec.ConstUnboundInt, http.StatusForbidden),
		I.Const("StatusFound", qspec.ConstUnboundInt, http.StatusFound),
		I.Const("StatusGatewayTimeout", qspec.ConstUnboundInt, http.StatusGatewayTimeout),
		I.Const("StatusGone", qspec.ConstUnboundInt, http.StatusGone),
		I.Const("StatusHTTPVersionNotSupported", qspec.ConstUnboundInt, http.StatusHTTPVersionNotSupported),
		I.Const("StatusIMUsed", qspec.ConstUnboundInt, http.StatusIMUsed),
		I.Const("StatusInsufficientStorage", qspec.ConstUnboundInt, http.StatusInsufficientStorage),
		I.Const("StatusInternalServerError", qspec.ConstUnboundInt, http.StatusInternalServerError),
		I.Const("StatusLengthRequired", qspec.ConstUnboundInt, http.StatusLengthRequired),
		I.Const("StatusLocked", qspec.ConstUnboundInt, http.StatusLocked),
		I.Const("StatusLoopDetected", qspec.ConstUnboundInt, http.StatusLoopDetected),
		I.Const("StatusMethodNotAllowed", qspec.ConstUnboundInt, http.StatusMethodNotAllowed),
		I.Const("StatusMisdirectedRequest", qspec.ConstUnboundInt, http.StatusMisdirectedRequest),
		I.Const("StatusMovedPermanently", qspec.ConstUnboundInt, http.StatusMovedPermanently),
		I.Const("StatusMultiStatus", qspec.ConstUnboundInt, http.StatusMultiStatus),
		I.Const("StatusMultipleChoices", qspec.ConstUnboundInt, http.StatusMultipleChoices),
		I.Const("StatusNetworkAuthenticationRequired", qspec.ConstUnboundInt, http.StatusNetworkAuthenticationRequired),
		I.Const("StatusNoContent", qspec.ConstUnboundInt, http.StatusNoContent),
		I.Const("StatusNonAuthoritativeInfo", qspec.ConstUnboundInt, http.StatusNonAuthoritativeInfo),
		I.Const("StatusNotAcceptable", qspec.ConstUnboundInt, http.StatusNotAcceptable),
		I.Const("StatusNotExtended", qspec.ConstUnboundInt, http.StatusNotExtended),
		I.Const("StatusNotFound", qspec.ConstUnboundInt, http.StatusNotFound),
		I.Const("StatusNotImplemented", qspec.ConstUnboundInt, http.StatusNotImplemented),
		I.Const("StatusNotModified", qspec.ConstUnboundInt, http.StatusNotModified),
		I.Const("StatusOK", qspec.ConstUnboundInt, http.StatusOK),
		I.Const("StatusPartialContent", qspec.ConstUnboundInt, http.StatusPartialContent),
		I.Const("StatusPaymentRequired", qspec.ConstUnboundInt, http.StatusPaymentRequired),
		I.Const("StatusPermanentRedirect", qspec.ConstUnboundInt, http.StatusPermanentRedirect),
		I.Const("StatusPreconditionFailed", qspec.ConstUnboundInt, http.StatusPreconditionFailed),
		I.Const("StatusPreconditionRequired", qspec.ConstUnboundInt, http.StatusPreconditionRequired),
		I.Const("StatusProcessing", qspec.ConstUnboundInt, http.StatusProcessing),
		I.Const("StatusProxyAuthRequired", qspec.ConstUnboundInt, http.StatusProxyAuthRequired),
		I.Const("StatusRequestEntityTooLarge", qspec.ConstUnboundInt, http.StatusRequestEntityTooLarge),
		I.Const("StatusRequestHeaderFieldsTooLarge", qspec.ConstUnboundInt, http.StatusRequestHeaderFieldsTooLarge),
		I.Const("StatusRequestTimeout", qspec.ConstUnboundInt, http.StatusRequestTimeout),
		I.Const("StatusRequestURITooLong", qspec.ConstUnboundInt, http.StatusRequestURITooLong),
		I.Const("StatusRequestedRangeNotSatisfiable", qspec.ConstUnboundInt, http.StatusRequestedRangeNotSatisfiable),
		I.Const("StatusResetContent", qspec.ConstUnboundInt, http.StatusResetContent),
		I.Const("StatusSeeOther", qspec.ConstUnboundInt, http.StatusSeeOther),
		I.Const("StatusServiceUnavailable", qspec.ConstUnboundInt, http.StatusServiceUnavailable),
		I.Const("StatusSwitchingProtocols", qspec.ConstUnboundInt, http.StatusSwitchingProtocols),
		I.Const("StatusTeapot", qspec.ConstUnboundInt, http.StatusTeapot),
		I.Const("StatusTemporaryRedirect", qspec.ConstUnboundInt, http.StatusTemporaryRedirect),
		I.Const("StatusTooEarly", qspec.ConstUnboundInt, http.StatusTooEarly),
		I.Const("StatusTooManyRequests", qspec.ConstUnboundInt, http.StatusTooManyRequests),
		I.Const("StatusUnauthorized", qspec.ConstUnboundInt, http.StatusUnauthorized),
		I.Const("StatusUnavailableForLegalReasons", qspec.ConstUnboundInt, http.StatusUnavailableForLegalReasons),
		I.Const("StatusUnprocessableEntity", qspec.ConstUnboundInt, http.StatusUnprocessableEntity),
		I.Const("StatusUnsupportedMediaType", qspec.ConstUnboundInt, http.StatusUnsupportedMediaType),
		I.Const("StatusUpgradeRequired", qspec.ConstUnboundInt, http.StatusUpgradeRequired),
		I.Const("StatusUseProxy", qspec.ConstUnboundInt, http.StatusUseProxy),
		I.Const("StatusVariantAlsoNegotiates", qspec.ConstUnboundInt, http.StatusVariantAlsoNegotiates),
		I.Const("TimeFormat", qspec.ConstBoundString, http.TimeFormat),
		I.Const("TrailerPrefix", qspec.ConstBoundString, http.TrailerPrefix),
	)
}

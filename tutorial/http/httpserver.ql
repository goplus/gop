mux = http.serveMux()
mux.handle("/404", http.notFoundHandler())
mux.handleFunc("/", fn(w, req) {
	fprintln(w, "host:", req.host, "path:", req.URL)
})

err = http.listenAndServe(":8888", mux)
if err != nil {
	fprintln(os.stderr, err)
}

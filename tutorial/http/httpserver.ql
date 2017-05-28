mux = http.NewServeMux()
mux.Handle("/404", http.NotFoundHandler())
mux.HandleFunc("/", fn(w, req) {
	fprintln(w, "host:", req.Host, "path:", req.URL)
})

err = http.ListenAndServe(":8888", mux)
if err != nil {
	fprintln(os.Stderr, err)
}

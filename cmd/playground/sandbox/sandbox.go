package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const maxRunTime = 2 * time.Second

type Request struct {
	Body string `json:"body"`
}

type Response struct {
	Errors string  `json:"errors"`
	Events []Event `json:"events"`
}

func main() {
	http.HandleFunc("/compile", compileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("error decoding request: %v", err), http.StatusBadRequest)
		return
	}
	resp, err := compileAndRun(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, fmt.Sprintf("error encoding response: %v", err), http.StatusInternalServerError)
		return
	}
}

func compileAndRun(req *Request) (res *Response, err error) {
	tmpDir, err := ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, fmt.Errorf("error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	in := filepath.Join(tmpDir, "main.qn")
	if err := ioutil.WriteFile(in, []byte(req.Body), 0400); err != nil {
		return nil, fmt.Errorf("error creating temp file %q: %v", in, err)
	}

	cmd := exec.Command("qlang.safe", in)
	rec := new(Recorder)
	cmd.Stdout = rec.Stdout()
	cmd.Stderr = rec.Stderr()

	if err := runTimeout(cmd, maxRunTime); err != nil {
		if err == timeoutErr {
			return &Response{Errors: "process took too long"}, nil
		}
		if _, ok := err.(*exec.ExitError); !ok {
			return nil, fmt.Errorf("error running sandbox: %v", err)
		}
	}
	events, err := rec.Events()
	if err != nil {
		return nil, fmt.Errorf("error decoding events: %v", err)
	}
	return &Response{Events: events}, nil
}

var timeoutErr = errors.New("process timed out")

func runTimeout(cmd *exec.Cmd, d time.Duration) error {
	if err := cmd.Start(); err != nil {
		return err
	}
	errc := make(chan error, 1)
	go func() {
		errc <- cmd.Wait()
	}()
	t := time.NewTimer(d)
	select {
	case err := <-errc:
		t.Stop()
		return err
	case <-t.C:
		cmd.Process.Kill()
		return timeoutErr
	}
}

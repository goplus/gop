// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package help implements the ``go help'' command.
package help

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/qiniu/x/log"
)

// Help implements the 'help' command.
func Help(w io.Writer, args []string) {
	cmd := base.Gop
Args:
	for i, arg := range args {
		for _, sub := range cmd.Commands {
			if sub.Name() == arg {
				cmd = sub
				continue Args
			}
		}

		// helpSuccess is the help command using as many args as possible that would succeed.
		helpSuccess := "gop help"
		if i > 0 {
			helpSuccess += " " + strings.Join(args[:i], " ")
		}
		fmt.Fprintf(os.Stderr, "gop help %s: unknown help topic. Run '%s'.\n", strings.Join(args, " "), helpSuccess)
		os.Exit(2)
	}

	if len(cmd.Commands) > 0 {
		PrintUsage(w, cmd)
	} else {
		cmd.Usage(w)
	}
	// not exit 2: succeeded at 'gop help cmd'.
	return
}

var usageTemplate = `{{.Short | trim}}

Usage:

	{{.UsageLine}} <command> [arguments]

The commands are:
{{range .Commands}}{{if or (.Runnable) .Commands}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "gop help{{with .LongName}} {{.}}{{end}} <command>" for more information about a command.

`

// commentWriter writes a Go comment to the underlying io.Writer,
// using line comment form (//).
type commentWriter struct {
	W            io.Writer
	wroteSlashes bool // Wrote "//" at the beginning of the current line.
}

func (c *commentWriter) Write(p []byte) (int, error) {
	var n int
	for i, b := range p {
		if !c.wroteSlashes {
			s := "//"
			if b != '\n' {
				s = "// "
			}
			if _, err := io.WriteString(c.W, s); err != nil {
				return n, err
			}
			c.wroteSlashes = true
		}
		n0, err := c.W.Write(p[i : i+1])
		n += n0
		if err != nil {
			return n, err
		}
		if b == '\n' {
			c.wroteSlashes = false
		}
	}
	return len(p), nil
}

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") {
			os.Exit(1)
		}
		log.Fatalf("writing output: %v", ew.err)
	}
	if err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

// PrintUsage prints usage information.
func PrintUsage(w io.Writer, cmd *base.Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, cmd)
	bw.Flush()
}

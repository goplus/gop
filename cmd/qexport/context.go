package main

import (
	"go/build"
	"log"
	"strings"
)

// contexts are the default contexts which are scanned, unless
// overridden by the -contexts flag.
var contexts = []*build.Context{
	{GOOS: "linux", GOARCH: "386", CgoEnabled: true},
	{GOOS: "linux", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "linux", GOARCH: "arm", CgoEnabled: true},
	{GOOS: "darwin", GOARCH: "386", CgoEnabled: true},
	{GOOS: "darwin", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "windows", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "windows", GOARCH: "386", CgoEnabled: true},
	{GOOS: "freebsd", GOARCH: "386", CgoEnabled: true},
	{GOOS: "freebsd", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "freebsd", GOARCH: "arm", CgoEnabled: true},
	{GOOS: "netbsd", GOARCH: "386", CgoEnabled: true},
	{GOOS: "netbsd", GOARCH: "amd64", CgoEnabled: true},
	{GOOS: "netbsd", GOARCH: "arm", CgoEnabled: true},
	{GOOS: "openbsd", GOARCH: "386", CgoEnabled: true},
	{GOOS: "openbsd", GOARCH: "amd64", CgoEnabled: true},
}

func contextName(c *build.Context) string {
	s := c.GOOS + "-" + c.GOARCH
	if c.CgoEnabled {
		return s + "-cgo"
	}
	return s
}

func osArchName(c *build.Context) string {
	return c.GOOS + "-" + c.GOARCH
}

func parseContext(c string) *build.Context {
	parts := strings.Split(c, "-")
	if len(parts) < 2 {
		log.Fatalf("bad context: %q", c)
	}
	bc := &build.Context{
		GOOS:   parts[0],
		GOARCH: parts[1],
	}
	if len(parts) == 3 {
		if parts[2] == "cgo" {
			bc.CgoEnabled = true
		} else {
			log.Fatalf("bad context: %q", c)
		}
	}
	return bc
}

func setCustomContexts(customctx string) {
	contexts = []*build.Context{}
	for _, c := range strings.Split(customctx, ",") {
		contexts = append(contexts, parseContext(c))
	}
}

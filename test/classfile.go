package test

import (
	"os"
	"testing"
)

const (
	GopPackage = true
)

// -----------------------------------------------------------------------------

type testingT = testing.T

// Case represents a Go+ testcase.
type Case struct {
	t *testingT
}

func (p *Case) initCase(t *testing.T) {
	p.t = t
}

// T returns a *testing.T object.
func (p Case) T() *testing.T { return p.t }

// Run runs f as a subtest of t called name. It runs f in a separate goroutine
// and blocks until f returns or calls t.Parallel to become a parallel test.
// Run reports whether f succeeded (or at least did not fail before calling t.Parallel).
func (p Case) Run(name string, f func(t *testing.T)) bool {
	return p.t.Run(name, f)
}

// Gopt_Case_TestMain is required by Go+ compiler as the test case entry.
func Gopt_Case_TestMain(c interface{ initCase(t *testing.T) }, t *testing.T) {
	c.initCase(t)
	c.(interface{ Main() }).Main()
}

// -----------------------------------------------------------------------------

// App represents a Go+ testing main application.
type App struct {
	m *testing.M
}

func (p *App) initApp(m *testing.M) {
	p.m = m
}

// Gopt_App_TestMain is required by Go+ compiler as the entry of a Go+ testing project.
func Gopt_App_TestMain(app interface{ initApp(m *testing.M) }, m *testing.M) {
	app.initApp(m)
	if me, ok := app.(interface{ MainEntry() }); ok {
		me.MainEntry()
	}
	os.Exit(m.Run())
}

// -----------------------------------------------------------------------------

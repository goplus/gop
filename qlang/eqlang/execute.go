package eqlang

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// -----------------------------------------------------------------------------

type qlang interface {
	Variables

	// SetVar sets a variable value.
	SetVar(name string, v interface{})

	// ResetVars resets all variables in executing context.
	ResetVars(vars map[string]interface{})

	// SafeExec compiles and executes a source code, without panic (will convert panic into an error).
	SafeExec(code []byte, fname string) (err error)
}

// -----------------------------------------------------------------------------

// A Type is a eql script engine.
//
type Type struct {
	Impl qlang
}

// New creates a new eql script engine.
//
func New(ql qlang) Type {

	return Type{ql}
}

// Var returns a variable value or defval if not found.
//
func (p Type) Var(name string, defval interface{}) interface{} {

	v, ok := p.Impl.GetVar(name)
	if !ok {
		v = defval
	}
	return v
}

// Imports returns import table.
//
func (p Type) Imports() string {

	imports := p.Var("imports", "").(string)
	if imports == "" {
		return ""
	}
	mods := strings.Split(imports, ",")
	return "\"" + strings.Join(mods, "\"\n\t\"") + "\""
}

// Subst substs variables in text.
//
func (p Type) Subst(text string) string {

	return subst(text, p.Impl)
}

// ExecuteDir executes a eql template directory.
//
func (p Type) ExecuteDir(global map[string]interface{}, source, output string) (err error) {

	if output == "" {
		output = p.Subst(source)
		if output == source {
			panic(fmt.Sprintf("source `%s` doesn't have $var", source))
		}
	}

	err = os.MkdirAll(output, 0755)
	if err != nil {
		return
	}

	fis, err := ioutil.ReadDir(source)
	if err != nil {
		return
	}

	source += "/"
	output += "/"
	for _, fi := range fis {
		name := fi.Name()
		if fi.IsDir() {
			err = p.ExecuteDir(global, source+name, output+name)
		} else if path.Ext(name) == ".eql" {
			p.Impl.ResetVars(global)
			newname := name[:len(name)-4]
			err = p.ExecuteFile(source+name, output+newname)
		} else {
			err = copyFile(source+name, output+name, fi.Mode())
		}
		if err != nil {
			return
		}
	}
	return
}

// ExecuteFile executes a eql template file.
//
func (p Type) ExecuteFile(source, output string) (err error) {

	b, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}

	return p.Execute(string(b), source, output)
}

// Execute executes a eql template string.
//
func (p Type) Execute(source string, fname string, output string) (err error) {

	p.Impl.SetVar("eql", p)
	code, err := Parse(source)
	if err != nil {
		return
	}

	if output != "" {
		f, err1 := os.Create(output)
		if err1 != nil {
			return err1
		}
		old := os.Stdout
		os.Stdout = f
		defer func() {
			os.Stdout = old
			f.Close()
		}()
	}

	err = p.Impl.SafeExec(code, fname)
	if err != nil && output != "" {
		os.Remove(output)
	}
	return
}

func copyFile(source, output string, perm os.FileMode) (err error) {

	b, err := ioutil.ReadFile(source)
	if err != nil {
		return
	}
	return ioutil.WriteFile(output, b, perm)
}

// -----------------------------------------------------------------------------

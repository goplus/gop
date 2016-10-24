package qlang

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"qiniupkg.com/text/tpl.v1/interpreter"
	"qlang.io/exec"
)

// -----------------------------------------------------------------------------

// An FindEntryError represents a FindEntry error.
//
type FindEntryError struct {
	Name string
	Err  error
}

func (e *FindEntryError) Error() string {
	return strconv.Quote(e.Name) + ": " + e.Err.Error()
}

func findEntry(file string, libs []string) (string, error) {

	if strings.HasPrefix(file, "/") {
		_, err := os.Stat(file)
		if err == nil {
			return file, nil
		}
		return "", &FindEntryError{file, err}
	}
	for _, dir := range libs {
		if dir == "" {
			continue
		}
		path := dir + "/" + file
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", &FindEntryError{file, syscall.ENOENT}
}

func resolvePath(file string, base string) string {

	if strings.HasPrefix(file, "/") {
		return file
	}
	return path.Join(base, file)
}

var (
	// FindEntry specifies the policy how qlang searches library file.
	FindEntry = findEntry

	// ReadFile specifies the policy how qlang reads source file.
	ReadFile = ioutil.ReadFile
)

// -----------------------------------------------------------------------------

const (
	indexFile = "/main.ql"
)

func qlangFile(file string) string {

	if path.Ext(file) == ".ql" {
		return file
	}
	return file + indexFile
}

func (p *Compiler) dir() string {

	if v, ok := p.gvars["__dir__"]; ok {
		if dir, ok := v.(string); ok {
			return dir
		}
	}
	panic("ident `__dir__` not found")
}

func (p *Compiler) compileModule(fname string, parent *funcCtx) (end int, symtbl map[string]int) {

	ctx := newFuncCtx(parent, nil)
	old := p.fnctx
	p.fnctx = ctx
	end = p.Compile(fname)
	symtbl = ctx.symtbl
	p.fnctx = old
	return
}

// Compile compiles a qlang source file.
//
func (p *Compiler) Compile(fname string) int {

	codeText, err := ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return p.Cl(codeText, fname)
}

// Cl compiles a qlang source code.
//
func (p *Compiler) Cl(codeText []byte, fname string) int {

	engine, err := interpreter.New(p, p.Opts)
	if err != nil {
		panic(err)
	}

	old := p.ipt
	defer func() { // backup interpreter to restore
		p.ipt = old
	}()
	p.ipt = engine
	p.gvars["__dir__"] = path.Dir(fname)
	p.gvars["__file__"] = fname
	err = engine.MatchExactly(codeText, fname)
	if err != nil {
		panic(err)
	}
	return p.code.Len()
}

// -----------------------------------------------------------------------------

func (p *Compiler) include(lit string) {

	file, err := strconv.Unquote(lit)
	if err != nil {
		panic("invalid string `" + lit + "`: " + err.Error())
	}

	fname := qlangFile(resolvePath(file, p.dir()))
	p.Compile(fname)
}

func requireModuleSym(fnctx *funcCtx, name string) int {

	if _, ok := fnctx.getSymbol(name); ok {
		panic("import `" + name + "` error: ident exists")
	}
	return fnctx.newSymbol(name)
}

func (p *Compiler) fnImport(lit string) {

	dir, err := strconv.Unquote(lit)
	if err != nil {
		panic("invalid string `" + lit + "`: " + err.Error())
	}

	file, err := FindEntry(dir+indexFile, p.libs)
	if err != nil {
		panic(err)
	}

	var name string
	arity := p.popArity()
	if arity > 0 {
		name = p.popName()
	} else {
		name = path.Base(dir)
		if pos := strings.Index(name, "."); pos > 0 {
			name = name[:pos]
		}
	}

	code := p.code
	instr := code.Reserve()
	id := requireModuleSym(p.fnctx, name)
	code.Block(exec.As(id))
	p.exits = append(p.exits, func() {
		file = path.Clean(file)
		mod, ok := p.mods[file]
		if !ok {
			start := code.Len()
			end, symtbl := p.compileModule(file, nil)
			mod = module{start, end, symtbl}
			p.mods[file] = mod
		}
		instr.Set(exec.Module(file, mod.start, mod.end, mod.symtbl))
	})
}

func (p *Compiler) export() {

	arity := p.popArity()
	names := p.gstk.PopFnArgs(arity)
	p.code.Block(exec.Export(names...))
}

// -----------------------------------------------------------------------------

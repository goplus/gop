/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package doc

import (
	"fmt"
	"go/types"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"

	"github.com/goplus/gop"
	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/cl/outline"
	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/x/gopenv"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gox"
)

// -----------------------------------------------------------------------------

// gop doc
var Cmd = &base.Command{
	UsageLine: "gop doc [-u -all -debug] [pkgPath]",
	Short:     "Show documentation for package or symbol",
}

var (
	flag    = &Cmd.Flag
	withDoc = flag.Bool("all", false, "Show all the documentation for the package.")
	debug   = flag.Bool("debug", false, "Print debug information.")
	unexp   = flag.Bool("u", false, "Show documentation for unexported as well as exported symbols, methods, and fields.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(cmd *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}

	pattern := flag.Args()
	if len(pattern) == 0 {
		pattern = []string{"."}
	}

	proj, _, err := gopprojs.ParseOne(pattern...)
	if err != nil {
		log.Panicln("gopprojs.ParseOne:", err)
	}

	if *debug {
		gox.SetDebug(gox.DbgFlagAll &^ gox.DbgFlagComments)
		cl.SetDebug(cl.DbgFlagAll)
		cl.SetDisableRecover(true)
	}

	gopEnv := gopenv.Get()
	conf := &gop.Config{Gop: gopEnv}
	outlinePkg(proj, conf)
}

func outlinePkg(proj gopprojs.Proj, conf *gop.Config) {
	var obj string
	var out outline.Package
	var err error
	switch v := proj.(type) {
	case *gopprojs.DirProj:
		obj = v.Dir
		out, err = gop.Outline(obj, conf)
	case *gopprojs.PkgPathProj:
		obj = v.Path
		out, err = gop.OutlinePkgPath("", obj, conf, true)
	default:
		log.Panicln("`gop doc` doesn't support", reflect.TypeOf(v))
	}
	if err == syscall.ENOENT {
		fmt.Fprintf(os.Stderr, "gop doc %v: not Go/Go+ files found\n", obj)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		outlineDoc(out.Pkg, out.Outline(*unexp), *unexp, *withDoc)
	}
}

const (
	indent = "    "
	ln     = "\n"
)

func outlineDoc(pkg *types.Package, out *outline.All, all, withDoc bool) {
	fmt.Printf("package %s // import %s\n\n", pkg.Name(), strconv.Quote(pkg.Path()))
	if withDoc && len(out.Consts) > 0 {
		fmt.Print("CONSTANTS\n\n")
	}
	for _, o := range out.Consts {
		fmt.Print(objectString(pkg, o.Const), ln)
	}
	if withDoc && len(out.Vars) > 0 {
		fmt.Print("VARIABLES\n\n")
	}
	for _, o := range out.Vars {
		fmt.Print(objectString(pkg, o.Var), ln)
	}
	if withDoc && len(out.Funcs) > 0 {
		fmt.Print("FUNCTIONS\n\n")
	}
	for _, fn := range out.Funcs {
		fmt.Print(objectString(pkg, fn.Obj()), ln)
		if withDoc {
			printDoc(fn)
		}
	}
	if withDoc && len(out.Types) > 0 {
		fmt.Print("TYPES\n\n")
	}
	for _, t := range out.Types {
		fmt.Println(typeString(pkg, t.TypeName, withDoc))
		for _, o := range t.Consts {
			fmt.Print(indent, constShortString(o.Const), ln)
		}
		for _, fn := range t.Creators {
			fmt.Print(indent, objectString(pkg, fn.Obj()), ln)
		}
		typ := t.Type()
		if named, ok := typ.CheckNamed(); ok {
			for _, fn := range named.Methods() {
				if o := fn.Obj(); all || o.Exported() {
					fmt.Print(indent, objectString(pkg, o), ln)
				}
			}
		}
	}
}

func printDoc(o interface{ Doc() string }) {
	fmt.Print(indent, strings.ReplaceAll(o.Doc(), "\n", "\n"+indent), ln)
}

func typeString(pkg *types.Package, t *types.TypeName, withDoc bool) string {
	alias := ""
	if t.IsAlias() {
		alias = " ="
	}
	underlying := typeShortString(pkg, t.Type().Underlying())
	return fmt.Sprint("type ", t.Name(), alias, " ", underlying)
}

func objectString(pkg *types.Package, obj types.Object) string {
	return types.ObjectString(obj, qualifier(pkg))
}

func typeShortString(pkg *types.Package, typ types.Type) string {
	switch t := typ.(type) {
	default:
		return types.TypeString(t, qualifier(pkg))
	}
}

func constShortString(obj *types.Const) string {
	return "const " + obj.Name()
}

func qualifier(pkg *types.Package) types.Qualifier {
	return func(other *types.Package) string {
		if pkg == other {
			return "" // same package; unqualified
		}
		return other.Name()
	}
}

// -----------------------------------------------------------------------------

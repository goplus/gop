package typesutil_test

import (
	"fmt"

	goast "go/ast"
	"go/constant"
	goformat "go/format"
	goparser "go/parser"

	"go/importer"
	"go/types"
	"sort"
	"strings"
	"testing"
	"unsafe"

	"github.com/goplus/gop"
	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/typesutil"
	"github.com/goplus/mod/env"
	"github.com/goplus/mod/gopmod"
)

func parserMixedSource(fset *token.FileSet, src string, gosrc string) (*typesutil.Info, *types.Info, error) {
	f, err := parser.ParseEntry(fset, "main.gop", src, parser.Config{
		Mode: parser.ParseComments,
	})
	if err != nil {
		return nil, nil, err
	}
	var gofiles []*goast.File
	if len(gosrc) > 0 {
		f, err := goparser.ParseFile(fset, "main.go", gosrc, goparser.ParseComments)
		if err != nil {
			return nil, nil, err
		}
		gofiles = append(gofiles, f)
	}

	conf := &types.Config{}
	conf.Importer = gop.NewImporter(nil, &env.Gop{Root: "../..", Version: "1.0"}, fset)
	chkOpts := &typesutil.Config{
		Types: types.NewPackage("main", "main"),
		Fset:  fset,
		Mod:   gopmod.Default,
	}
	info := &typesutil.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	ginfo := &types.Info{
		Types:      make(map[goast.Expr]types.TypeAndValue),
		Defs:       make(map[*goast.Ident]types.Object),
		Uses:       make(map[*goast.Ident]types.Object),
		Implicits:  make(map[goast.Node]types.Object),
		Selections: make(map[*goast.SelectorExpr]*types.Selection),
		Scopes:     make(map[goast.Node]*types.Scope),
	}
	check := typesutil.NewChecker(conf, chkOpts, ginfo, info)
	err = check.Files(gofiles, []*ast.File{f})
	return info, ginfo, err
}

func parserSource(fset *token.FileSet, filename string, src interface{}, mode parser.Mode) (*typesutil.Info, error) {
	f, err := parser.ParseEntry(fset, filename, src, parser.Config{
		Mode: mode,
	})
	if err != nil {
		return nil, err
	}

	conf := &types.Config{}
	conf.Importer = importer.Default()
	chkOpts := &typesutil.Config{
		Types: types.NewPackage("main", "main"),
		Fset:  fset,
		Mod:   gopmod.Default,
	}
	info := &typesutil.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	check := typesutil.NewChecker(conf, chkOpts, nil, info)
	err = check.Files(nil, []*ast.File{f})
	return info, err
}

func parserGoSource(fset *token.FileSet, filename string, src interface{}, mode goparser.Mode) (*types.Info, error) {
	f, err := goparser.ParseFile(fset, filename, src, mode)
	if err != nil {
		return nil, err
	}

	conf := &types.Config{}
	conf.Importer = importer.Default()
	info := &types.Info{
		Types:      make(map[goast.Expr]types.TypeAndValue),
		Defs:       make(map[*goast.Ident]types.Object),
		Uses:       make(map[*goast.Ident]types.Object),
		Implicits:  make(map[goast.Node]types.Object),
		Selections: make(map[*goast.SelectorExpr]*types.Selection),
		Scopes:     make(map[goast.Node]*types.Scope),
	}
	pkg := types.NewPackage("main", "main")
	check := types.NewChecker(conf, fset, pkg, info)
	err = check.Files([]*goast.File{f})
	return info, err
}

func testGopInfo(t *testing.T, src string, gosrc string, expect string) {
	fset := token.NewFileSet()
	info, _, err := parserMixedSource(fset, src, gosrc)
	if err != nil {
		t.Fatal("parserMixedSource error", err)
	}
	var list []string
	list = append(list, "== types ==")
	list = append(list, typesList(fset, info.Types, false)...)
	list = append(list, "== defs ==")
	list = append(list, defsList(fset, info.Defs, true)...)
	list = append(list, "== uses ==")
	list = append(list, usesList(fset, info.Uses)...)
	result := strings.Join(list, "\n")
	t.Log(result)
	if result != expect {
		t.Fatal("bad expect\n", expect)
	}
}

func testInfo(t *testing.T, src interface{}) {
	fset := token.NewFileSet()
	info, err := parserSource(fset, "main.gop", src, parser.ParseComments)
	if err != nil {
		t.Fatal("parserSource error", err)
	}
	goinfo, err := parserGoSource(fset, "main.go", src, goparser.ParseComments)
	if err != nil {
		t.Fatal("parserGoSource error", err)
	}
	testItems(t, "types", typesList(fset, info.Types, true), goTypesList(fset, goinfo.Types, true))
	testItems(t, "defs", defsList(fset, info.Defs, true), goDefsList(fset, goinfo.Defs, true))
	testItems(t, "uses", usesList(fset, info.Uses), goUsesList(fset, goinfo.Uses))
	// TODO check selections
	//testItems(t, "selections", selectionList(fset, info.Selections), goSelectionList(fset, goinfo.Selections))
}

func testItems(t *testing.T, name string, items []string, goitems []string) {
	text := strings.Join(items, "\n")
	gotext := strings.Join(goitems, "\n")
	if len(items) != len(goitems) || text != gotext {
		t.Errorf(`====== check %v error (Go+ count: %v, Go count %v) ====== 
------ Go+ ------
%v
------ Go ------
%v
`,
			name, len(items), len(goitems),
			text, gotext)
	} else {
		t.Log(fmt.Sprintf(`====== check %v pass (count: %v) ======
%v
`, name, len(items), text))
	}
}

func sortItems(items []string) []string {
	sort.Strings(items)
	for i := 0; i < len(items); i++ {
		items[i] = fmt.Sprintf("%03v: %v", i, items[i])
	}
	return items
}

func typesList(fset *token.FileSet, types map[ast.Expr]types.TypeAndValue, skipBasicLit bool) []string {
	var items []string
	for expr, tv := range types {
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		tvstr := tv.Type.String()
		if skipBasicLit {
			if t, ok := expr.(*ast.BasicLit); ok {
				tvstr = t.Kind.String()
			}
		}
		if tv.Value != nil {
			tvstr += " = " + tv.Value.String()
		}
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s %-30T | %-7s : %s | %v",
			posn.Line, posn.Column, exprString(fset, expr), expr,
			mode(tv), tvstr, (*TypeAndValue)(unsafe.Pointer(&tv)).mode)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func goTypesList(fset *token.FileSet, types map[goast.Expr]types.TypeAndValue, skipBasicLit bool) []string {
	var items []string
	for expr, tv := range types {
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		tvstr := tv.Type.String()
		if skipBasicLit {
			if t, ok := expr.(*goast.BasicLit); ok {
				tvstr = t.Kind.String()
			}
		}
		if tv.Value != nil {
			tvstr += " = " + tv.Value.String()
		}
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s %-30T | %-7s : %s | %v",
			posn.Line, posn.Column, goexprString(fset, expr), expr,
			mode(tv), tvstr, (*TypeAndValue)(unsafe.Pointer(&tv)).mode)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func defsList(fset *token.FileSet, uses map[*ast.Ident]types.Object, skipNil bool) []string {
	var items []string
	for expr, obj := range uses {
		if skipNil && obj == nil {
			continue
		}
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s | %s",
			posn.Line, posn.Column, expr,
			obj)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func goDefsList(fset *token.FileSet, uses map[*goast.Ident]types.Object, skipNil bool) []string {
	var items []string
	for expr, obj := range uses {
		if skipNil && obj == nil {
			continue // skip nil object
		}
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s | %s",
			posn.Line, posn.Column, expr,
			obj)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func usesList(fset *token.FileSet, uses map[*ast.Ident]types.Object) []string {
	var items []string
	for expr, obj := range uses {
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s | %s",
			posn.Line, posn.Column, expr,
			obj)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func goUsesList(fset *token.FileSet, uses map[*goast.Ident]types.Object) []string {
	var items []string
	for expr, obj := range uses {
		if obj == nil {
			continue // skip nil object
		}
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s | %s",
			posn.Line, posn.Column, expr,
			obj)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func selectionList(fset *token.FileSet, sels map[*ast.SelectorExpr]*types.Selection) []string {
	var items []string
	for expr, sel := range sels {
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s | %s",
			posn.Line, posn.Column, exprString(fset, expr),
			sel)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func goSelectionList(fset *token.FileSet, sels map[*goast.SelectorExpr]*types.Selection) []string {
	var items []string
	for expr, sel := range sels {
		var buf strings.Builder
		posn := fset.Position(expr.Pos())
		// line:col | expr | mode : type = value
		fmt.Fprintf(&buf, "%2d:%2d | %-19s | %s",
			posn.Line, posn.Column, goexprString(fset, expr),
			sel)
		items = append(items, buf.String())
	}
	return sortItems(items)
}

func mode(tv types.TypeAndValue) string {
	switch {
	case tv.IsVoid():
		return "void"
	case tv.IsType():
		return "type"
	case tv.IsBuiltin():
		return "builtin"
	case tv.IsNil():
		return "nil"
	case tv.Assignable():
		if tv.Addressable() {
			return "var"
		}
		return "mapindex"
	case tv.IsValue():
		return "value"
	default:
		return "unknown"
	}
}

func exprString(fset *token.FileSet, expr ast.Expr) string {
	var buf strings.Builder
	format.Node(&buf, fset, expr)
	return buf.String()
}

func goexprString(fset *token.FileSet, expr goast.Expr) string {
	var buf strings.Builder
	goformat.Node(&buf, fset, expr)
	return buf.String()
}

type operandMode byte

const (
	invalid   operandMode = iota // operand is invalid
	novalue                      // operand represents no value (result of a function call w/o result)
	builtin                      // operand is a built-in function
	typexpr                      // operand is a type
	constant_                    // operand is a constant; the operand's typ is a Basic type
	variable                     // operand is an addressable variable
	mapindex                     // operand is a map index expression (acts like a variable on lhs, commaok on rhs of an assignment)
	value                        // operand is a computed value
	commaok                      // like value, but operand may be used in a comma,ok expression
	commaerr                     // like commaok, but second value is error, not boolean
	cgofunc                      // operand is a cgo function
)

func (v operandMode) String() string {
	return operandModeString[int(v)]
}

var operandModeString = [...]string{
	invalid:   "invalid operand",
	novalue:   "no value",
	builtin:   "built-in",
	typexpr:   "type",
	constant_: "constant",
	variable:  "variable",
	mapindex:  "map index expression",
	value:     "value",
	commaok:   "comma, ok expression",
	commaerr:  "comma, error expression",
	cgofunc:   "cgo function",
}

type TypeAndValue struct {
	mode  operandMode
	Type  types.Type
	Value constant.Value
}

func TestVarTypes(t *testing.T) {
	testInfo(t, `package main
type T struct {
	x int
	y int
}
var v *int = nil
var v1 []int
var v2 map[int8]string
var v3 struct{}
var v4 *T = &T{100,200}
var v5 = [6]int{}
var v6 = v5[0]
var v7 = [6]int{}[0]
var m map[int]string
func init() {
	v5[0] = 100
	_ = v5[:][0]
	m[0] = "hello"
	_ = m[0]
	_ = map[int]string{}[0]
	_ = &v3
	_ = *(&v3)
	a := []int{1,2,3,4,5}[0]
	_ = a
}
`)
}

func TestStruct(t *testing.T) {
	testInfo(t, `package main
import "fmt"

type Person struct {
	name string
	age  int8
}

func test() {
	p := Person{
		name: "jack",
	}
	_ = p.name
	p.name = "name"
	fmt.Println(p)
}
`)
}

func TestTypeAssert(t *testing.T) {
	testInfo(t, `package main

func test() {
	var a interface{} = 100
	if n, ok := a.(int); ok {
		_ = n
	}
}
`)
}

func TestChan(t *testing.T) {
	testInfo(t, `package main

func test() {
	var ch chan int
	select {
	case n, ok := <-ch:
		_ = n
		_ = ok
		break
	}
}
`)
}

func TestRange(t *testing.T) {
	testInfo(t, `package main
func test() {
	a := []int{100,200}
	for k, v := range a {
		_ = k
		_ = v
	}
	var m map[int]string
	for k, v := range m {
		_ = k
		_ = v
	}
	for v := range m {
		_ = v
	}
}
`)
}

func TestFuncLit(t *testing.T) {
	testInfo(t, `package main
func test() {
	add := func(n1 int, n2 int) int {
		return n1+n2
	}
	_ = add(1,2)

	go func(n int) {
		_ = n+100
	}(100)
}
`)
}

func TestSliceLit(t *testing.T) {
	testGopInfo(t, `
a := [100,200]
println a
`, ``, `== types ==
000:  2: 7 | 100                 *ast.BasicLit                  | value   : untyped int = 100 | constant
001:  2:11 | 200                 *ast.BasicLit                  | value   : untyped int = 200 | constant
002:  3: 1 | println             *ast.Ident                     | value   : func(a ...any) (n int, err error) | value
003:  3: 1 | println a           *ast.CallExpr                  | value   : (n int, err error) | value
004:  3: 9 | a                   *ast.Ident                     | var     : []int | variable
== defs ==
000:  2: 1 | a                   | var a []int
001:  2: 1 | main                | func main.main()
== uses ==
000:  3: 1 | println             | func fmt.Println(a ...any) (n int, err error)
001:  3: 9 | a                   | var a []int`)
}

func TestForPhrase1(t *testing.T) {
	testGopInfo(t, `
sum := 0
for x <- [1, 3, 5, 7, 11, 13, 17], x > 3 {
	sum = sum + x
}
println sum
`, ``, `== types ==
000:  2: 8 | 0                   *ast.BasicLit                  | value   : untyped int = 0 | constant
001:  3:11 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
002:  3:14 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
003:  3:17 | 5                   *ast.BasicLit                  | value   : untyped int = 5 | constant
004:  3:20 | 7                   *ast.BasicLit                  | value   : untyped int = 7 | constant
005:  3:23 | 11                  *ast.BasicLit                  | value   : untyped int = 11 | constant
006:  3:27 | 13                  *ast.BasicLit                  | value   : untyped int = 13 | constant
007:  3:31 | 17                  *ast.BasicLit                  | value   : untyped int = 17 | constant
008:  3:36 | x                   *ast.Ident                     | var     : int | variable
009:  3:36 | x > 3               *ast.BinaryExpr                | value   : untyped bool | value
010:  3:40 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
011:  4: 2 | sum                 *ast.Ident                     | var     : int | variable
012:  4: 8 | sum                 *ast.Ident                     | var     : int | variable
013:  4: 8 | sum + x             *ast.BinaryExpr                | value   : int | value
014:  4:14 | x                   *ast.Ident                     | var     : int | variable
015:  6: 1 | println             *ast.Ident                     | value   : func(a ...any) (n int, err error) | value
016:  6: 1 | println sum         *ast.CallExpr                  | value   : (n int, err error) | value
017:  6: 9 | sum                 *ast.Ident                     | var     : int | variable
== defs ==
000:  2: 1 | main                | func main.main()
001:  2: 1 | sum                 | var sum int
002:  3: 5 | x                   | var x int
== uses ==
000:  3:36 | x                   | var x int
001:  4: 2 | sum                 | var sum int
002:  4: 8 | sum                 | var sum int
003:  4:14 | x                   | var x int
004:  6: 1 | println             | func fmt.Println(a ...any) (n int, err error)
005:  6: 9 | sum                 | var sum int`)
}

func TestForPhrase2(t *testing.T) {
	testGopInfo(t, `
sum := 0
for i, x <- [1, 3, 5, 7, 11, 13, 17], i%2 == 1 && x > 3 {
	sum = sum + x
}
println sum
`, ``, `== types ==
000:  2: 8 | 0                   *ast.BasicLit                  | value   : untyped int = 0 | constant
001:  3:14 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
002:  3:17 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
003:  3:20 | 5                   *ast.BasicLit                  | value   : untyped int = 5 | constant
004:  3:23 | 7                   *ast.BasicLit                  | value   : untyped int = 7 | constant
005:  3:26 | 11                  *ast.BasicLit                  | value   : untyped int = 11 | constant
006:  3:30 | 13                  *ast.BasicLit                  | value   : untyped int = 13 | constant
007:  3:34 | 17                  *ast.BasicLit                  | value   : untyped int = 17 | constant
008:  3:39 | i                   *ast.Ident                     | var     : int | variable
009:  3:39 | i % 2               *ast.BinaryExpr                | value   : int | value
010:  3:39 | i%2 == 1            *ast.BinaryExpr                | value   : untyped bool | value
011:  3:39 | i%2 == 1 && x > 3   *ast.BinaryExpr                | value   : untyped bool | value
012:  3:41 | 2                   *ast.BasicLit                  | value   : untyped int = 2 | constant
013:  3:46 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
014:  3:51 | x                   *ast.Ident                     | var     : int | variable
015:  3:51 | x > 3               *ast.BinaryExpr                | value   : untyped bool | value
016:  3:55 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
017:  4: 2 | sum                 *ast.Ident                     | var     : int | variable
018:  4: 8 | sum                 *ast.Ident                     | var     : int | variable
019:  4: 8 | sum + x             *ast.BinaryExpr                | value   : int | value
020:  4:14 | x                   *ast.Ident                     | var     : int | variable
021:  6: 1 | println             *ast.Ident                     | value   : func(a ...any) (n int, err error) | value
022:  6: 1 | println sum         *ast.CallExpr                  | value   : (n int, err error) | value
023:  6: 9 | sum                 *ast.Ident                     | var     : int | variable
== defs ==
000:  2: 1 | main                | func main.main()
001:  2: 1 | sum                 | var sum int
002:  3: 5 | i                   | var i int
003:  3: 8 | x                   | var x int
== uses ==
000:  3:39 | i                   | var i int
001:  3:51 | x                   | var x int
002:  4: 2 | sum                 | var sum int
003:  4: 8 | sum                 | var sum int
004:  4:14 | x                   | var x int
005:  6: 1 | println             | func fmt.Println(a ...any) (n int, err error)
006:  6: 9 | sum                 | var sum int`)
}

func TestMapComprehension(t *testing.T) {
	testGopInfo(t, `
y := {x: i for i, x <- ["1", "3", "5", "7", "11"]}
println y
`, ``, `== types ==
000:  2: 7 | x                   *ast.Ident                     | var     : string | variable
001:  2:10 | i                   *ast.Ident                     | var     : int | variable
002:  2:25 | "1"                 *ast.BasicLit                  | value   : untyped string = "1" | constant
003:  2:30 | "3"                 *ast.BasicLit                  | value   : untyped string = "3" | constant
004:  2:35 | "5"                 *ast.BasicLit                  | value   : untyped string = "5" | constant
005:  2:40 | "7"                 *ast.BasicLit                  | value   : untyped string = "7" | constant
006:  2:45 | "11"                *ast.BasicLit                  | value   : untyped string = "11" | constant
007:  3: 1 | println             *ast.Ident                     | value   : func(a ...any) (n int, err error) | value
008:  3: 1 | println y           *ast.CallExpr                  | value   : (n int, err error) | value
009:  3: 9 | y                   *ast.Ident                     | var     : map[string]int | variable
== defs ==
000:  2: 1 | main                | func main.main()
001:  2: 1 | y                   | var y map[string]int
002:  2:16 | i                   | var i int
003:  2:19 | x                   | var x string
== uses ==
000:  2: 7 | x                   | var x string
001:  2:10 | i                   | var i int
002:  3: 1 | println             | func fmt.Println(a ...any) (n int, err error)
003:  3: 9 | y                   | var y map[string]int`)
}

func TestListComprehension(t *testing.T) {
	testGopInfo(t, `
a := [1, 3.4, 5]
b := [x*x for x <- a]
_ = b
`, ``, `== types ==
000:  2: 7 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
001:  2:10 | 3.4                 *ast.BasicLit                  | value   : untyped float = 3.4 | constant
002:  2:15 | 5                   *ast.BasicLit                  | value   : untyped int = 5 | constant
003:  3: 7 | x                   *ast.Ident                     | var     : float64 | variable
004:  3: 7 | x * x               *ast.BinaryExpr                | value   : float64 | value
005:  3: 9 | x                   *ast.Ident                     | var     : float64 | variable
006:  3:20 | a                   *ast.Ident                     | var     : []float64 | variable
007:  4: 5 | b                   *ast.Ident                     | var     : []float64 | variable
== defs ==
000:  2: 1 | a                   | var a []float64
001:  2: 1 | main                | func main.main()
002:  3: 1 | b                   | var b []float64
003:  3:15 | x                   | var x float64
== uses ==
000:  3: 7 | x                   | var x float64
001:  3: 9 | x                   | var x float64
002:  3:20 | a                   | var a []float64
003:  4: 5 | b                   | var b []float64`)
}

func TestListComprehensionMultiLevel(t *testing.T) {
	testGopInfo(t, `
arr := [1, 2, 3, 4.1, 5, 6]
x := [[a, b] for a <- arr, a < b for b <- arr, b > 2]
println("x:", x)
`, ``, `== types ==
000:  2: 9 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
001:  2:12 | 2                   *ast.BasicLit                  | value   : untyped int = 2 | constant
002:  2:15 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
003:  2:18 | 4.1                 *ast.BasicLit                  | value   : untyped float = 4.1 | constant
004:  2:23 | 5                   *ast.BasicLit                  | value   : untyped int = 5 | constant
005:  2:26 | 6                   *ast.BasicLit                  | value   : untyped int = 6 | constant
006:  3: 8 | a                   *ast.Ident                     | var     : float64 | variable
007:  3:11 | b                   *ast.Ident                     | var     : float64 | variable
008:  3:23 | arr                 *ast.Ident                     | var     : []float64 | variable
009:  3:28 | a                   *ast.Ident                     | var     : float64 | variable
010:  3:28 | a < b               *ast.BinaryExpr                | value   : untyped bool | value
011:  3:32 | b                   *ast.Ident                     | var     : float64 | variable
012:  3:43 | arr                 *ast.Ident                     | var     : []float64 | variable
013:  3:48 | b                   *ast.Ident                     | var     : float64 | variable
014:  3:48 | b > 2               *ast.BinaryExpr                | value   : untyped bool | value
015:  3:52 | 2                   *ast.BasicLit                  | value   : untyped int = 2 | constant
016:  4: 1 | println             *ast.Ident                     | value   : func(a ...any) (n int, err error) | value
017:  4: 1 | println("x:", x)    *ast.CallExpr                  | value   : (n int, err error) | value
018:  4: 9 | "x:"                *ast.BasicLit                  | value   : untyped string = "x:" | constant
019:  4:15 | x                   *ast.Ident                     | var     : [][]float64 | variable
== defs ==
000:  2: 1 | arr                 | var arr []float64
001:  2: 1 | main                | func main.main()
002:  3: 1 | x                   | var x [][]float64
003:  3:18 | a                   | var a float64
004:  3:38 | b                   | var b float64
== uses ==
000:  3: 8 | a                   | var a float64
001:  3:11 | b                   | var b float64
002:  3:23 | arr                 | var arr []float64
003:  3:28 | a                   | var a float64
004:  3:32 | b                   | var b float64
005:  3:43 | arr                 | var arr []float64
006:  3:48 | b                   | var b float64
007:  4: 1 | println             | func fmt.Println(a ...any) (n int, err error)
008:  4:15 | x                   | var x [][]float64`)
}

func TestFileEnumLines(t *testing.T) {
	testGopInfo(t, `
import "os"

for line <- os.Stdin {
	println line
}
`, ``, `== types ==
000:  4:13 | os.Stdin            *ast.SelectorExpr              | var     : *os.File | variable
001:  5: 2 | println             *ast.Ident                     | value   : func(a ...any) (n int, err error) | value
002:  5: 2 | println line        *ast.CallExpr                  | value   : (n int, err error) | value
003:  5:10 | line                *ast.Ident                     | var     : string | variable
== defs ==
000:  4: 1 | main                | func main.main()
001:  4: 5 | line                | var line string
== uses ==
000:  4:13 | os                  | package os
001:  4:16 | Stdin               | var os.Stdin *os.File
002:  5: 2 | println             | func fmt.Println(a ...any) (n int, err error)
003:  5:10 | line                | var line string`)
}

func TestLambdaExpr(t *testing.T) {
	testGopInfo(t, `package main
func Map(c []float64, t func(float64) float64) {
	// ...
}

func Map2(c []float64, t func(float64) (float64, float64)) {
	// ...
}

Map([1.2, 3.5, 6], x => x * x)
Map2([1.2, 3.5, 6], x => (x * x, x + x))
`, ``, `== types ==
000:  2:12 | []float64           *ast.ArrayType                 | type    : []float64 | type
001:  2:14 | float64             *ast.Ident                     | type    : float64 | type
002:  2:25 | func(float64) float64 *ast.FuncType                  | type    : func(float64) float64 | type
003:  2:30 | float64             *ast.Ident                     | type    : float64 | type
004:  2:39 | float64             *ast.Ident                     | type    : float64 | type
005:  6:13 | []float64           *ast.ArrayType                 | type    : []float64 | type
006:  6:15 | float64             *ast.Ident                     | type    : float64 | type
007:  6:26 | func(float64) (float64, float64) *ast.FuncType                  | type    : func(float64) (float64, float64) | type
008:  6:31 | float64             *ast.Ident                     | type    : float64 | type
009:  6:41 | float64             *ast.Ident                     | type    : float64 | type
010:  6:50 | float64             *ast.Ident                     | type    : float64 | type
011: 10: 1 | Map                 *ast.Ident                     | value   : func(c []float64, t func(float64) float64) | value
012: 10: 1 | Map([1.2, 3.5, 6], x => x * x) *ast.CallExpr                  | void    : () | no value
013: 10: 6 | 1.2                 *ast.BasicLit                  | value   : untyped float = 1.2 | constant
014: 10:11 | 3.5                 *ast.BasicLit                  | value   : untyped float = 3.5 | constant
015: 10:16 | 6                   *ast.BasicLit                  | value   : untyped int = 6 | constant
016: 10:25 | x                   *ast.Ident                     | var     : float64 | variable
017: 10:25 | x * x               *ast.BinaryExpr                | value   : float64 | value
018: 10:29 | x                   *ast.Ident                     | var     : float64 | variable
019: 11: 1 | Map2                *ast.Ident                     | value   : func(c []float64, t func(float64) (float64, float64)) | value
020: 11: 1 | Map2([1.2, 3.5, 6], x => (x * x, x + x)) *ast.CallExpr                  | void    : () | no value
021: 11: 7 | 1.2                 *ast.BasicLit                  | value   : untyped float = 1.2 | constant
022: 11:12 | 3.5                 *ast.BasicLit                  | value   : untyped float = 3.5 | constant
023: 11:17 | 6                   *ast.BasicLit                  | value   : untyped int = 6 | constant
024: 11:27 | x                   *ast.Ident                     | var     : float64 | variable
025: 11:27 | x * x               *ast.BinaryExpr                | value   : float64 | value
026: 11:31 | x                   *ast.Ident                     | var     : float64 | variable
027: 11:34 | x                   *ast.Ident                     | var     : float64 | variable
028: 11:34 | x + x               *ast.BinaryExpr                | value   : float64 | value
029: 11:38 | x                   *ast.Ident                     | var     : float64 | variable
== defs ==
000:  2: 6 | Map                 | func main.Map(c []float64, t func(float64) float64)
001:  2:10 | c                   | var c []float64
002:  2:23 | t                   | var t func(float64) float64
003:  6: 6 | Map2                | func main.Map2(c []float64, t func(float64) (float64, float64))
004:  6:11 | c                   | var c []float64
005:  6:24 | t                   | var t func(float64) (float64, float64)
006: 10: 1 | main                | func main.main()
007: 10:20 | x                   | var x float64
008: 11:21 | x                   | var x float64
== uses ==
000:  2:14 | float64             | type float64
001:  2:30 | float64             | type float64
002:  2:39 | float64             | type float64
003:  6:15 | float64             | type float64
004:  6:31 | float64             | type float64
005:  6:41 | float64             | type float64
006:  6:50 | float64             | type float64
007: 10: 1 | Map                 | func main.Map(c []float64, t func(float64) float64)
008: 10:25 | x                   | var x float64
009: 10:29 | x                   | var x float64
010: 11: 1 | Map2                | func main.Map2(c []float64, t func(float64) (float64, float64))
011: 11:27 | x                   | var x float64
012: 11:31 | x                   | var x float64
013: 11:34 | x                   | var x float64
014: 11:38 | x                   | var x float64`)
}

func TestLambdaExpr2(t *testing.T) {
	testGopInfo(t, `package main
func Map(c []float64, t func(float64) float64) {
	// ...
}

func Map2(c []float64, t func(float64) (float64, float64)) {
	// ...
}

Map([1.2, 3.5, 6], x => {
	return x * x
})
Map2([1.2, 3.5, 6], x => {
	return x * x, x + x
})
`, ``, `== types ==
000:  2:12 | []float64           *ast.ArrayType                 | type    : []float64 | type
001:  2:14 | float64             *ast.Ident                     | type    : float64 | type
002:  2:25 | func(float64) float64 *ast.FuncType                  | type    : func(float64) float64 | type
003:  2:30 | float64             *ast.Ident                     | type    : float64 | type
004:  2:39 | float64             *ast.Ident                     | type    : float64 | type
005:  6:13 | []float64           *ast.ArrayType                 | type    : []float64 | type
006:  6:15 | float64             *ast.Ident                     | type    : float64 | type
007:  6:26 | func(float64) (float64, float64) *ast.FuncType                  | type    : func(float64) (float64, float64) | type
008:  6:31 | float64             *ast.Ident                     | type    : float64 | type
009:  6:41 | float64             *ast.Ident                     | type    : float64 | type
010:  6:50 | float64             *ast.Ident                     | type    : float64 | type
011: 10: 1 | Map                 *ast.Ident                     | value   : func(c []float64, t func(float64) float64) | value
012: 10: 1 | Map([1.2, 3.5, 6], x => {
	return x * x
}) *ast.CallExpr                  | void    : () | no value
013: 10: 6 | 1.2                 *ast.BasicLit                  | value   : untyped float = 1.2 | constant
014: 10:11 | 3.5                 *ast.BasicLit                  | value   : untyped float = 3.5 | constant
015: 10:16 | 6                   *ast.BasicLit                  | value   : untyped int = 6 | constant
016: 11: 9 | x                   *ast.Ident                     | var     : float64 | variable
017: 11: 9 | x * x               *ast.BinaryExpr                | value   : float64 | value
018: 11:13 | x                   *ast.Ident                     | var     : float64 | variable
019: 13: 1 | Map2                *ast.Ident                     | value   : func(c []float64, t func(float64) (float64, float64)) | value
020: 13: 1 | Map2([1.2, 3.5, 6], x => {
	return x * x, x + x
}) *ast.CallExpr                  | void    : () | no value
021: 13: 7 | 1.2                 *ast.BasicLit                  | value   : untyped float = 1.2 | constant
022: 13:12 | 3.5                 *ast.BasicLit                  | value   : untyped float = 3.5 | constant
023: 13:17 | 6                   *ast.BasicLit                  | value   : untyped int = 6 | constant
024: 14: 9 | x                   *ast.Ident                     | var     : float64 | variable
025: 14: 9 | x * x               *ast.BinaryExpr                | value   : float64 | value
026: 14:13 | x                   *ast.Ident                     | var     : float64 | variable
027: 14:16 | x                   *ast.Ident                     | var     : float64 | variable
028: 14:16 | x + x               *ast.BinaryExpr                | value   : float64 | value
029: 14:20 | x                   *ast.Ident                     | var     : float64 | variable
== defs ==
000:  2: 6 | Map                 | func main.Map(c []float64, t func(float64) float64)
001:  2:10 | c                   | var c []float64
002:  2:23 | t                   | var t func(float64) float64
003:  6: 6 | Map2                | func main.Map2(c []float64, t func(float64) (float64, float64))
004:  6:11 | c                   | var c []float64
005:  6:24 | t                   | var t func(float64) (float64, float64)
006: 10: 1 | main                | func main.main()
007: 10:20 | x                   | var x float64
008: 13:21 | x                   | var x float64
== uses ==
000:  2:14 | float64             | type float64
001:  2:30 | float64             | type float64
002:  2:39 | float64             | type float64
003:  6:15 | float64             | type float64
004:  6:31 | float64             | type float64
005:  6:41 | float64             | type float64
006:  6:50 | float64             | type float64
007: 10: 1 | Map                 | func main.Map(c []float64, t func(float64) float64)
008: 11: 9 | x                   | var x float64
009: 11:13 | x                   | var x float64
010: 13: 1 | Map2                | func main.Map2(c []float64, t func(float64) (float64, float64))
011: 14: 9 | x                   | var x float64
012: 14:13 | x                   | var x float64
013: 14:16 | x                   | var x float64
014: 14:20 | x                   | var x float64`)
}

func TestMixedOverload1(t *testing.T) {
	testGopInfo(t, `
type Mesh struct {
}

func (p *Mesh) Name() string {
	return "hello"
}

var (
	m1 = &Mesh{}
	m2 = &Mesh{}
)

OnKey "hello", => {
}
OnKey "hello", key => {
}
OnKey ["1"], => {
}
OnKey ["2"], key => {
}
OnKey [m1, m2], => {
}
OnKey [m1, m2], key => {
}
OnKey ["a"], ["b"], key => {
}
OnKey ["a"], [m1, m2], key => {
}
OnKey ["a"], nil, key => {
}
OnKey 100, 200
OnKey "a", "b", x => x * x, x => {
	return x * 2
}
OnKey "a", "b", 1, 2, 3
OnKey("a", "b", [1, 2, 3]...)
`, `
package main

type Mesher interface {
	Name() string
}

type N struct {
}

func (m *N) OnKey__0(a string, fn func()) {
}

func (m *N) OnKey__1(a string, fn func(key string)) {
}

func (m *N) OnKey__2(a []string, fn func()) {
}

func (m *N) OnKey__3(a []string, fn func(key string)) {
}

func (m *N) OnKey__4(a []Mesher, fn func()) {
}

func (m *N) OnKey__5(a []Mesher, fn func(key Mesher)) {
}

func (m *N) OnKey__6(a []string, b []string, fn func(key string)) {
}

func (m *N) OnKey__7(a []string, b []Mesher, fn func(key string)) {
}

func (m *N) OnKey__8(x int, y int) {
}


func OnKey__0(a string, fn func()) {
}

func OnKey__1(a string, fn func(key string)) {
}

func OnKey__2(a []string, fn func()) {
}

func OnKey__3(a []string, fn func(key string)) {
}

func OnKey__4(a []Mesher, fn func()) {
}

func OnKey__5(a []Mesher, fn func(key Mesher)) {
}

func OnKey__6(a []string, b []string, fn func(key string)) {
}

func OnKey__7(a []string, b []Mesher, fn func(key string)) {
}

func OnKey__8(x int, y int) {
}

func OnKey__9(a, b string, fn ...func(x int) int) {
}

func OnKey__a(a, b string, v ...int) {
}
`, `== types ==
000:  2:11 | struct {
}          *ast.StructType                | type    : struct{} | type
001:  5:10 | Mesh                *ast.Ident                     | type    : main.Mesh | type
002:  5:23 | string              *ast.Ident                     | type    : string | type
003:  6: 9 | "hello"             *ast.BasicLit                  | value   : untyped string = "hello" | constant
004: 10: 7 | &Mesh{}             *ast.UnaryExpr                 | value   : *main.Mesh | value
005: 10: 8 | Mesh                *ast.Ident                     | type    : main.Mesh | type
006: 10: 8 | Mesh{}              *ast.CompositeLit              | value   : main.Mesh | value
007: 11: 7 | &Mesh{}             *ast.UnaryExpr                 | value   : *main.Mesh | value
008: 11: 8 | Mesh                *ast.Ident                     | type    : main.Mesh | type
009: 11: 8 | Mesh{}              *ast.CompositeLit              | value   : main.Mesh | value
010: 14: 1 | OnKey               *ast.Ident                     | value   : func(a string, fn func()) | value
011: 14: 1 | OnKey "hello", => {
} *ast.CallExpr                  | void    : () | no value
012: 14: 7 | "hello"             *ast.BasicLit                  | value   : untyped string = "hello" | constant
013: 16: 1 | OnKey               *ast.Ident                     | value   : func(a string, fn func(key string)) | value
014: 16: 1 | OnKey "hello", key => {
} *ast.CallExpr                  | void    : () | no value
015: 16: 7 | "hello"             *ast.BasicLit                  | value   : untyped string = "hello" | constant
016: 18: 1 | OnKey               *ast.Ident                     | value   : func(a []string, fn func()) | value
017: 18: 1 | OnKey ["1"], => {
} *ast.CallExpr                  | void    : () | no value
018: 18: 8 | "1"                 *ast.BasicLit                  | value   : untyped string = "1" | constant
019: 20: 1 | OnKey               *ast.Ident                     | value   : func(a []string, fn func(key string)) | value
020: 20: 1 | OnKey ["2"], key => {
} *ast.CallExpr                  | void    : () | no value
021: 20: 8 | "2"                 *ast.BasicLit                  | value   : untyped string = "2" | constant
022: 22: 1 | OnKey               *ast.Ident                     | value   : func(a []main.Mesher, fn func()) | value
023: 22: 1 | OnKey [m1, m2], => {
} *ast.CallExpr                  | void    : () | no value
024: 22: 8 | m1                  *ast.Ident                     | var     : *main.Mesh | variable
025: 22:12 | m2                  *ast.Ident                     | var     : *main.Mesh | variable
026: 24: 1 | OnKey               *ast.Ident                     | value   : func(a []main.Mesher, fn func(key main.Mesher)) | value
027: 24: 1 | OnKey [m1, m2], key => {
} *ast.CallExpr                  | void    : () | no value
028: 24: 8 | m1                  *ast.Ident                     | var     : *main.Mesh | variable
029: 24:12 | m2                  *ast.Ident                     | var     : *main.Mesh | variable
030: 26: 1 | OnKey               *ast.Ident                     | value   : func(a []string, b []string, fn func(key string)) | value
031: 26: 1 | OnKey ["a"], ["b"], key => {
} *ast.CallExpr                  | void    : () | no value
032: 26: 8 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
033: 26:15 | "b"                 *ast.BasicLit                  | value   : untyped string = "b" | constant
034: 28: 1 | OnKey               *ast.Ident                     | value   : func(a []string, b []main.Mesher, fn func(key string)) | value
035: 28: 1 | OnKey ["a"], [m1, m2], key => {
} *ast.CallExpr                  | void    : () | no value
036: 28: 8 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
037: 28:15 | m1                  *ast.Ident                     | var     : *main.Mesh | variable
038: 28:19 | m2                  *ast.Ident                     | var     : *main.Mesh | variable
039: 30: 1 | OnKey               *ast.Ident                     | value   : func(a []string, b []string, fn func(key string)) | value
040: 30: 1 | OnKey ["a"], nil, key => {
} *ast.CallExpr                  | void    : () | no value
041: 30: 8 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
042: 30:14 | nil                 *ast.Ident                     | nil     : untyped nil | value
043: 32: 1 | OnKey               *ast.Ident                     | value   : func(x int, y int) | value
044: 32: 1 | OnKey 100, 200      *ast.CallExpr                  | void    : () | no value
045: 32: 7 | 100                 *ast.BasicLit                  | value   : untyped int = 100 | constant
046: 32:12 | 200                 *ast.BasicLit                  | value   : untyped int = 200 | constant
047: 33: 1 | OnKey               *ast.Ident                     | value   : func(a string, b string, fn ...func(x int) int) | value
048: 33: 1 | OnKey "a", "b", x => x * x, x => {
	return x * 2
} *ast.CallExpr                  | void    : () | no value
049: 33: 7 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
050: 33:12 | "b"                 *ast.BasicLit                  | value   : untyped string = "b" | constant
051: 33:22 | x                   *ast.Ident                     | var     : int | variable
052: 33:22 | x * x               *ast.BinaryExpr                | value   : int | value
053: 33:26 | x                   *ast.Ident                     | var     : int | variable
054: 34: 9 | x                   *ast.Ident                     | var     : int | variable
055: 34: 9 | x * 2               *ast.BinaryExpr                | value   : int | value
056: 34:13 | 2                   *ast.BasicLit                  | value   : untyped int = 2 | constant
057: 36: 1 | OnKey               *ast.Ident                     | value   : func(a string, b string, v ...int) | value
058: 36: 1 | OnKey "a", "b", 1, 2, 3 *ast.CallExpr                  | void    : () | no value
059: 36: 7 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
060: 36:12 | "b"                 *ast.BasicLit                  | value   : untyped string = "b" | constant
061: 36:17 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
062: 36:20 | 2                   *ast.BasicLit                  | value   : untyped int = 2 | constant
063: 36:23 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
064: 37: 1 | OnKey               *ast.Ident                     | value   : func(a string, b string, v ...int) | value
065: 37: 1 | OnKey("a", "b", [1, 2, 3]...) *ast.CallExpr                  | void    : () | no value
066: 37: 7 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
067: 37:12 | "b"                 *ast.BasicLit                  | value   : untyped string = "b" | constant
068: 37:18 | 1                   *ast.BasicLit                  | value   : untyped int = 1 | constant
069: 37:21 | 2                   *ast.BasicLit                  | value   : untyped int = 2 | constant
070: 37:24 | 3                   *ast.BasicLit                  | value   : untyped int = 3 | constant
== defs ==
000:  2: 6 | Mesh                | type main.Mesh struct{}
001:  5: 7 | p                   | var p *main.Mesh
002:  5:16 | Name                | func (*main.Mesh).Name() string
003: 10: 2 | m1                  | var main.m1 *main.Mesh
004: 11: 2 | m2                  | var main.m2 *main.Mesh
005: 14: 1 | main                | func main.main()
006: 16:16 | key                 | var key string
007: 20:14 | key                 | var key string
008: 24:17 | key                 | var key main.Mesher
009: 26:21 | key                 | var key string
010: 28:24 | key                 | var key string
011: 30:19 | key                 | var key string
012: 33:17 | x                   | var x int
013: 33:29 | x                   | var x int
== uses ==
000:  5:10 | Mesh                | type main.Mesh struct{}
001:  5:23 | string              | type string
002: 10: 8 | Mesh                | type main.Mesh struct{}
003: 11: 8 | Mesh                | type main.Mesh struct{}
004: 14: 1 | OnKey               | func main.OnKey__0(a string, fn func())
005: 16: 1 | OnKey               | func main.OnKey__1(a string, fn func(key string))
006: 18: 1 | OnKey               | func main.OnKey__2(a []string, fn func())
007: 20: 1 | OnKey               | func main.OnKey__3(a []string, fn func(key string))
008: 22: 1 | OnKey               | func main.OnKey__4(a []main.Mesher, fn func())
009: 22: 8 | m1                  | var main.m1 *main.Mesh
010: 22:12 | m2                  | var main.m2 *main.Mesh
011: 24: 1 | OnKey               | func main.OnKey__5(a []main.Mesher, fn func(key main.Mesher))
012: 24: 8 | m1                  | var main.m1 *main.Mesh
013: 24:12 | m2                  | var main.m2 *main.Mesh
014: 26: 1 | OnKey               | func main.OnKey__6(a []string, b []string, fn func(key string))
015: 28: 1 | OnKey               | func main.OnKey__7(a []string, b []main.Mesher, fn func(key string))
016: 28:15 | m1                  | var main.m1 *main.Mesh
017: 28:19 | m2                  | var main.m2 *main.Mesh
018: 30: 1 | OnKey               | func main.OnKey__6(a []string, b []string, fn func(key string))
019: 30:14 | nil                 | nil
020: 32: 1 | OnKey               | func main.OnKey__8(x int, y int)
021: 33: 1 | OnKey               | func main.OnKey__9(a string, b string, fn ...func(x int) int)
022: 33:22 | x                   | var x int
023: 33:26 | x                   | var x int
024: 34: 9 | x                   | var x int
025: 36: 1 | OnKey               | func main.OnKey__a(a string, b string, v ...int)
026: 37: 1 | OnKey               | func main.OnKey__a(a string, b string, v ...int)`)
}

func TestMixedOverload2(t *testing.T) {
	testGopInfo(t, `
type Mesh struct {
}

func (p *Mesh) Name() string {
	return "hello"
}

var (
	m1 = &Mesh{}
	m2 = &Mesh{}
)

n := &N{}
n.onKey "hello", => {
}
n.onKey "hello", key => {
}
n.onKey ["1"], => {
}
n.onKey ["2"], key => {
}
n.onKey [m1, m2], => {
}
n.onKey [m1, m2], key => {
}
n.onKey ["a"], ["b"], key => {
}
n.onKey ["a"], [m1, m2], key => {
}
n.onKey ["a"], nil, key => {
}
n.onKey 100, 200
`, `
package main

type Mesher interface {
	Name() string
}

type N struct {
}

func (m *N) OnKey__0(a string, fn func()) {
}

func (m *N) OnKey__1(a string, fn func(key string)) {
}

func (m *N) OnKey__2(a []string, fn func()) {
}

func (m *N) OnKey__3(a []string, fn func(key string)) {
}

func (m *N) OnKey__4(a []Mesher, fn func()) {
}

func (m *N) OnKey__5(a []Mesher, fn func(key Mesher)) {
}

func (m *N) OnKey__6(a []string, b []string, fn func(key string)) {
}

func (m *N) OnKey__7(a []string, b []Mesher, fn func(key string)) {
}

func (m *N) OnKey__8(x int, y int) {
}


func OnKey__0(a string, fn func()) {
}

func OnKey__1(a string, fn func(key string)) {
}

func OnKey__2(a []string, fn func()) {
}

func OnKey__3(a []string, fn func(key string)) {
}

func OnKey__4(a []Mesher, fn func()) {
}

func OnKey__5(a []Mesher, fn func(key Mesher)) {
}

func OnKey__6(a []string, b []string, fn func(key string)) {
}

func OnKey__7(a []string, b []Mesher, fn func(key string)) {
}

func OnKey__8(x int, y int) {
}

func OnKey__9(a, b string, fn ...func(x int) int) {
}

func OnKey__a(a, b string, v ...int) {
}
`, `== types ==
000:  2:11 | struct {
}          *ast.StructType                | type    : struct{} | type
001:  5:10 | Mesh                *ast.Ident                     | type    : main.Mesh | type
002:  5:23 | string              *ast.Ident                     | type    : string | type
003:  6: 9 | "hello"             *ast.BasicLit                  | value   : untyped string = "hello" | constant
004: 10: 7 | &Mesh{}             *ast.UnaryExpr                 | value   : *main.Mesh | value
005: 10: 8 | Mesh                *ast.Ident                     | type    : main.Mesh | type
006: 10: 8 | Mesh{}              *ast.CompositeLit              | value   : main.Mesh | value
007: 11: 7 | &Mesh{}             *ast.UnaryExpr                 | value   : *main.Mesh | value
008: 11: 8 | Mesh                *ast.Ident                     | type    : main.Mesh | type
009: 11: 8 | Mesh{}              *ast.CompositeLit              | value   : main.Mesh | value
010: 14: 6 | &N{}                *ast.UnaryExpr                 | value   : *main.N | value
011: 14: 7 | N                   *ast.Ident                     | type    : main.N | type
012: 14: 7 | N{}                 *ast.CompositeLit              | value   : main.N | value
013: 15: 1 | n                   *ast.Ident                     | var     : *main.N | variable
014: 15: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a string, fn func()) | value
015: 15: 1 | n.onKey "hello", => {
} *ast.CallExpr                  | void    : () | no value
016: 15: 9 | "hello"             *ast.BasicLit                  | value   : untyped string = "hello" | constant
017: 17: 1 | n                   *ast.Ident                     | var     : *main.N | variable
018: 17: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a string, fn func(key string)) | value
019: 17: 1 | n.onKey "hello", key => {
} *ast.CallExpr                  | void    : () | no value
020: 17: 9 | "hello"             *ast.BasicLit                  | value   : untyped string = "hello" | constant
021: 19: 1 | n                   *ast.Ident                     | var     : *main.N | variable
022: 19: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []string, fn func()) | value
023: 19: 1 | n.onKey ["1"], => {
} *ast.CallExpr                  | void    : () | no value
024: 19:10 | "1"                 *ast.BasicLit                  | value   : untyped string = "1" | constant
025: 21: 1 | n                   *ast.Ident                     | var     : *main.N | variable
026: 21: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []string, fn func(key string)) | value
027: 21: 1 | n.onKey ["2"], key => {
} *ast.CallExpr                  | void    : () | no value
028: 21:10 | "2"                 *ast.BasicLit                  | value   : untyped string = "2" | constant
029: 23: 1 | n                   *ast.Ident                     | var     : *main.N | variable
030: 23: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []main.Mesher, fn func()) | value
031: 23: 1 | n.onKey [m1, m2], => {
} *ast.CallExpr                  | void    : () | no value
032: 23:10 | m1                  *ast.Ident                     | var     : *main.Mesh | variable
033: 23:14 | m2                  *ast.Ident                     | var     : *main.Mesh | variable
034: 25: 1 | n                   *ast.Ident                     | var     : *main.N | variable
035: 25: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []main.Mesher, fn func(key main.Mesher)) | value
036: 25: 1 | n.onKey [m1, m2], key => {
} *ast.CallExpr                  | void    : () | no value
037: 25:10 | m1                  *ast.Ident                     | var     : *main.Mesh | variable
038: 25:14 | m2                  *ast.Ident                     | var     : *main.Mesh | variable
039: 27: 1 | n                   *ast.Ident                     | var     : *main.N | variable
040: 27: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []string, b []string, fn func(key string)) | value
041: 27: 1 | n.onKey ["a"], ["b"], key => {
} *ast.CallExpr                  | void    : () | no value
042: 27:10 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
043: 27:17 | "b"                 *ast.BasicLit                  | value   : untyped string = "b" | constant
044: 29: 1 | n                   *ast.Ident                     | var     : *main.N | variable
045: 29: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []string, b []main.Mesher, fn func(key string)) | value
046: 29: 1 | n.onKey ["a"], [m1, m2], key => {
} *ast.CallExpr                  | void    : () | no value
047: 29:10 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
048: 29:17 | m1                  *ast.Ident                     | var     : *main.Mesh | variable
049: 29:21 | m2                  *ast.Ident                     | var     : *main.Mesh | variable
050: 31: 1 | n                   *ast.Ident                     | var     : *main.N | variable
051: 31: 1 | n.onKey             *ast.SelectorExpr              | value   : func(a []string, b []string, fn func(key string)) | value
052: 31: 1 | n.onKey ["a"], nil, key => {
} *ast.CallExpr                  | void    : () | no value
053: 31:10 | "a"                 *ast.BasicLit                  | value   : untyped string = "a" | constant
054: 31:16 | nil                 *ast.Ident                     | nil     : untyped nil | value
055: 33: 1 | n                   *ast.Ident                     | var     : *main.N | variable
056: 33: 1 | n.onKey             *ast.SelectorExpr              | value   : func(x int, y int) | value
057: 33: 1 | n.onKey 100, 200    *ast.CallExpr                  | void    : () | no value
058: 33: 9 | 100                 *ast.BasicLit                  | value   : untyped int = 100 | constant
059: 33:14 | 200                 *ast.BasicLit                  | value   : untyped int = 200 | constant
== defs ==
000:  2: 6 | Mesh                | type main.Mesh struct{}
001:  5: 7 | p                   | var p *main.Mesh
002:  5:16 | Name                | func (*main.Mesh).Name() string
003: 10: 2 | m1                  | var main.m1 *main.Mesh
004: 11: 2 | m2                  | var main.m2 *main.Mesh
005: 14: 1 | main                | func main.main()
006: 14: 1 | n                   | var n *main.N
007: 17:18 | key                 | var key string
008: 21:16 | key                 | var key string
009: 25:19 | key                 | var key main.Mesher
010: 27:23 | key                 | var key string
011: 29:26 | key                 | var key string
012: 31:21 | key                 | var key string
== uses ==
000:  5:10 | Mesh                | type main.Mesh struct{}
001:  5:23 | string              | type string
002: 10: 8 | Mesh                | type main.Mesh struct{}
003: 11: 8 | Mesh                | type main.Mesh struct{}
004: 14: 7 | N                   | type main.N struct{}
005: 15: 1 | n                   | var n *main.N
006: 15: 3 | onKey               | func (*main.N).OnKey__0(a string, fn func())
007: 17: 1 | n                   | var n *main.N
008: 17: 3 | onKey               | func (*main.N).OnKey__1(a string, fn func(key string))
009: 19: 1 | n                   | var n *main.N
010: 19: 3 | onKey               | func (*main.N).OnKey__2(a []string, fn func())
011: 21: 1 | n                   | var n *main.N
012: 21: 3 | onKey               | func (*main.N).OnKey__3(a []string, fn func(key string))
013: 23: 1 | n                   | var n *main.N
014: 23: 3 | onKey               | func (*main.N).OnKey__4(a []main.Mesher, fn func())
015: 23:10 | m1                  | var main.m1 *main.Mesh
016: 23:14 | m2                  | var main.m2 *main.Mesh
017: 25: 1 | n                   | var n *main.N
018: 25: 3 | onKey               | func (*main.N).OnKey__5(a []main.Mesher, fn func(key main.Mesher))
019: 25:10 | m1                  | var main.m1 *main.Mesh
020: 25:14 | m2                  | var main.m2 *main.Mesh
021: 27: 1 | n                   | var n *main.N
022: 27: 3 | onKey               | func (*main.N).OnKey__6(a []string, b []string, fn func(key string))
023: 29: 1 | n                   | var n *main.N
024: 29: 3 | onKey               | func (*main.N).OnKey__7(a []string, b []main.Mesher, fn func(key string))
025: 29:17 | m1                  | var main.m1 *main.Mesh
026: 29:21 | m2                  | var main.m2 *main.Mesh
027: 31: 1 | n                   | var n *main.N
028: 31: 3 | onKey               | func (*main.N).OnKey__6(a []string, b []string, fn func(key string))
029: 31:16 | nil                 | nil
030: 33: 1 | n                   | var n *main.N
031: 33: 3 | onKey               | func (*main.N).OnKey__8(x int, y int)`)
}

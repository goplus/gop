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

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/format"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
	"github.com/goplus/gop/x/typesutil"
	"github.com/goplus/mod/gopmod"
)

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

func TestFuncLit(t *testing.T) {
	testInfo(t, `package main

func test() {
	add := func(n1 int, n2 int) int {
		return n1+n2
	}
	_ = add(1,2)
}
`)
}

func TestBuiltin(t *testing.T) {
	testInfo(t, `package main

func main() {
	var a = []int{100,200}
	a = append(a,100,200)
	_ = a
	var b []byte
	b = append(b,"world"...)
	_ = b
	_ = cap(a)
	_ = len(a)
	_ = new(int)
}
`)
}

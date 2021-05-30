// Package sort provide Go+ "sort" package, as "sort" package in Go.
package sort

import (
	reflect "reflect"
	sort "sort"

	gop "github.com/goplus/gop"
)

func execmFloat64SliceSearch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(sort.Float64Slice).Search(args[1].(float64))
	p.Ret(2, ret0)
}

func execmFloat64SliceLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(sort.Float64Slice).Len()
	p.Ret(1, ret0)
}

func execmFloat64SliceLess(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(sort.Float64Slice).Less(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmFloat64SliceSwap(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(sort.Float64Slice).Swap(args[1].(int), args[2].(int))
	p.PopN(3)
}

func execmFloat64SliceSort(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(sort.Float64Slice).Sort()
	p.PopN(1)
}

func execFloat64s(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	sort.Float64s(args[0].([]float64))
	p.PopN(1)
}

func execFloat64sAreSorted(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sort.Float64sAreSorted(args[0].([]float64))
	p.Ret(1, ret0)
}

func execmIntSliceSearch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(sort.IntSlice).Search(args[1].(int))
	p.Ret(2, ret0)
}

func execmIntSliceLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(sort.IntSlice).Len()
	p.Ret(1, ret0)
}

func execmIntSliceLess(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(sort.IntSlice).Less(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmIntSliceSwap(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(sort.IntSlice).Swap(args[1].(int), args[2].(int))
	p.PopN(3)
}

func execmIntSliceSort(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(sort.IntSlice).Sort()
	p.PopN(1)
}

func execiInterfaceLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(sort.Interface).Len()
	p.Ret(1, ret0)
}

func execiInterfaceLess(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(sort.Interface).Less(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execiInterfaceSwap(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(sort.Interface).Swap(args[1].(int), args[2].(int))
	p.PopN(3)
}

func execInts(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	sort.Ints(args[0].([]int))
	p.PopN(1)
}

func execIntsAreSorted(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sort.IntsAreSorted(args[0].([]int))
	p.Ret(1, ret0)
}

func toType0(v interface{}) sort.Interface {
	if v == nil {
		return nil
	}
	return v.(sort.Interface)
}

func execIsSorted(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sort.IsSorted(toType0(args[0]))
	p.Ret(1, ret0)
}

func execReverse(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sort.Reverse(toType0(args[0]))
	p.Ret(1, ret0)
}

func execSearch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := sort.Search(args[0].(int), args[1].(func(int) bool))
	p.Ret(2, ret0)
}

func execSearchFloat64s(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := sort.SearchFloat64s(args[0].([]float64), args[1].(float64))
	p.Ret(2, ret0)
}

func execSearchInts(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := sort.SearchInts(args[0].([]int), args[1].(int))
	p.Ret(2, ret0)
}

func execSearchStrings(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := sort.SearchStrings(args[0].([]string), args[1].(string))
	p.Ret(2, ret0)
}

func execSlice(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	sort.Slice(args[0], args[1].(func(i int, j int) bool))
	p.PopN(2)
}

func execSliceIsSorted(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := sort.SliceIsSorted(args[0], args[1].(func(i int, j int) bool))
	p.Ret(2, ret0)
}

func execSliceStable(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	sort.SliceStable(args[0], args[1].(func(i int, j int) bool))
	p.PopN(2)
}

func execSort(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	sort.Sort(toType0(args[0]))
	p.PopN(1)
}

func execStable(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	sort.Stable(toType0(args[0]))
	p.PopN(1)
}

func execmStringSliceSearch(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(sort.StringSlice).Search(args[1].(string))
	p.Ret(2, ret0)
}

func execmStringSliceLen(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(sort.StringSlice).Len()
	p.Ret(1, ret0)
}

func execmStringSliceLess(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(sort.StringSlice).Less(args[1].(int), args[2].(int))
	p.Ret(3, ret0)
}

func execmStringSliceSwap(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(sort.StringSlice).Swap(args[1].(int), args[2].(int))
	p.PopN(3)
}

func execmStringSliceSort(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(sort.StringSlice).Sort()
	p.PopN(1)
}

func execStrings(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	sort.Strings(args[0].([]string))
	p.PopN(1)
}

func execStringsAreSorted(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := sort.StringsAreSorted(args[0].([]string))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("sort")

func init() {
	I.RegisterFuncs(
		I.Func("(Float64Slice).Search", (sort.Float64Slice).Search, execmFloat64SliceSearch),
		I.Func("(Float64Slice).Len", (sort.Float64Slice).Len, execmFloat64SliceLen),
		I.Func("(Float64Slice).Less", (sort.Float64Slice).Less, execmFloat64SliceLess),
		I.Func("(Float64Slice).Swap", (sort.Float64Slice).Swap, execmFloat64SliceSwap),
		I.Func("(Float64Slice).Sort", (sort.Float64Slice).Sort, execmFloat64SliceSort),
		I.Func("Float64s", sort.Float64s, execFloat64s),
		I.Func("Float64sAreSorted", sort.Float64sAreSorted, execFloat64sAreSorted),
		I.Func("(IntSlice).Search", (sort.IntSlice).Search, execmIntSliceSearch),
		I.Func("(IntSlice).Len", (sort.IntSlice).Len, execmIntSliceLen),
		I.Func("(IntSlice).Less", (sort.IntSlice).Less, execmIntSliceLess),
		I.Func("(IntSlice).Swap", (sort.IntSlice).Swap, execmIntSliceSwap),
		I.Func("(IntSlice).Sort", (sort.IntSlice).Sort, execmIntSliceSort),
		I.Func("(Interface).Len", (sort.Interface).Len, execiInterfaceLen),
		I.Func("(Interface).Less", (sort.Interface).Less, execiInterfaceLess),
		I.Func("(Interface).Swap", (sort.Interface).Swap, execiInterfaceSwap),
		I.Func("Ints", sort.Ints, execInts),
		I.Func("IntsAreSorted", sort.IntsAreSorted, execIntsAreSorted),
		I.Func("IsSorted", sort.IsSorted, execIsSorted),
		I.Func("Reverse", sort.Reverse, execReverse),
		I.Func("Search", sort.Search, execSearch),
		I.Func("SearchFloat64s", sort.SearchFloat64s, execSearchFloat64s),
		I.Func("SearchInts", sort.SearchInts, execSearchInts),
		I.Func("SearchStrings", sort.SearchStrings, execSearchStrings),
		I.Func("Slice", sort.Slice, execSlice),
		I.Func("SliceIsSorted", sort.SliceIsSorted, execSliceIsSorted),
		I.Func("SliceStable", sort.SliceStable, execSliceStable),
		I.Func("Sort", sort.Sort, execSort),
		I.Func("Stable", sort.Stable, execStable),
		I.Func("(StringSlice).Search", (sort.StringSlice).Search, execmStringSliceSearch),
		I.Func("(StringSlice).Len", (sort.StringSlice).Len, execmStringSliceLen),
		I.Func("(StringSlice).Less", (sort.StringSlice).Less, execmStringSliceLess),
		I.Func("(StringSlice).Swap", (sort.StringSlice).Swap, execmStringSliceSwap),
		I.Func("(StringSlice).Sort", (sort.StringSlice).Sort, execmStringSliceSort),
		I.Func("Strings", sort.Strings, execStrings),
		I.Func("StringsAreSorted", sort.StringsAreSorted, execStringsAreSorted),
	)
	I.RegisterTypes(
		I.Type("Float64Slice", reflect.TypeOf((*sort.Float64Slice)(nil)).Elem()),
		I.Type("IntSlice", reflect.TypeOf((*sort.IntSlice)(nil)).Elem()),
		I.Type("Interface", reflect.TypeOf((*sort.Interface)(nil)).Elem()),
		I.Type("StringSlice", reflect.TypeOf((*sort.StringSlice)(nil)).Elem()),
	)
}

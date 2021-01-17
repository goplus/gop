// Package flag provide Go+ "flag" package, as "flag" package in Go.
package flag

import (
	flag "flag"
	io "io"
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execArg(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := flag.Arg(args[0].(int))
	p.Ret(1, ret0)
}

func execArgs(_ int, p *gop.Context) {
	ret0 := flag.Args()
	p.Ret(0, ret0)
}

func execBool(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Bool(args[0].(string), args[1].(bool), args[2].(string))
	p.Ret(3, ret0)
}

func execBoolVar(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.BoolVar(args[0].(*bool), args[1].(string), args[2].(bool), args[3].(string))
	p.PopN(4)
}

func execDuration(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Duration(args[0].(string), args[1].(time.Duration), args[2].(string))
	p.Ret(3, ret0)
}

func execDurationVar(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.DurationVar(args[0].(*time.Duration), args[1].(string), args[2].(time.Duration), args[3].(string))
	p.PopN(4)
}

func execmFlagSetOutput(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).Output()
	p.Ret(1, ret0)
}

func execmFlagSetName(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).Name()
	p.Ret(1, ret0)
}

func execmFlagSetErrorHandling(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).ErrorHandling()
	p.Ret(1, ret0)
}

func toType0(v interface{}) io.Writer {
	if v == nil {
		return nil
	}
	return v.(io.Writer)
}

func execmFlagSetSetOutput(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*flag.FlagSet).SetOutput(toType0(args[1]))
	p.PopN(2)
}

func execmFlagSetVisitAll(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*flag.FlagSet).VisitAll(args[1].(func(*flag.Flag)))
	p.PopN(2)
}

func execmFlagSetVisit(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	args[0].(*flag.FlagSet).Visit(args[1].(func(*flag.Flag)))
	p.PopN(2)
}

func execmFlagSetLookup(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*flag.FlagSet).Lookup(args[1].(string))
	p.Ret(2, ret0)
}

func execmFlagSetSet(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(*flag.FlagSet).Set(args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execmFlagSetPrintDefaults(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*flag.FlagSet).PrintDefaults()
	p.PopN(1)
}

func execmFlagSetNFlag(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).NFlag()
	p.Ret(1, ret0)
}

func execmFlagSetArg(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*flag.FlagSet).Arg(args[1].(int))
	p.Ret(2, ret0)
}

func execmFlagSetNArg(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).NArg()
	p.Ret(1, ret0)
}

func execmFlagSetArgs(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).Args()
	p.Ret(1, ret0)
}

func execmFlagSetBoolVar(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).BoolVar(args[1].(*bool), args[2].(string), args[3].(bool), args[4].(string))
	p.PopN(5)
}

func execmFlagSetBool(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Bool(args[1].(string), args[2].(bool), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetIntVar(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).IntVar(args[1].(*int), args[2].(string), args[3].(int), args[4].(string))
	p.PopN(5)
}

func execmFlagSetInt(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Int(args[1].(string), args[2].(int), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetInt64Var(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).Int64Var(args[1].(*int64), args[2].(string), args[3].(int64), args[4].(string))
	p.PopN(5)
}

func execmFlagSetInt64(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Int64(args[1].(string), args[2].(int64), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetUintVar(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).UintVar(args[1].(*uint), args[2].(string), args[3].(uint), args[4].(string))
	p.PopN(5)
}

func execmFlagSetUint(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Uint(args[1].(string), args[2].(uint), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetUint64Var(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).Uint64Var(args[1].(*uint64), args[2].(string), args[3].(uint64), args[4].(string))
	p.PopN(5)
}

func execmFlagSetUint64(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Uint64(args[1].(string), args[2].(uint64), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetStringVar(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).StringVar(args[1].(*string), args[2].(string), args[3].(string), args[4].(string))
	p.PopN(5)
}

func execmFlagSetString(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).String(args[1].(string), args[2].(string), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetFloat64Var(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).Float64Var(args[1].(*float64), args[2].(string), args[3].(float64), args[4].(string))
	p.PopN(5)
}

func execmFlagSetFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Float64(args[1].(string), args[2].(float64), args[3].(string))
	p.Ret(4, ret0)
}

func execmFlagSetDurationVar(_ int, p *gop.Context) {
	args := p.GetArgs(5)
	args[0].(*flag.FlagSet).DurationVar(args[1].(*time.Duration), args[2].(string), args[3].(time.Duration), args[4].(string))
	p.PopN(5)
}

func execmFlagSetDuration(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(*flag.FlagSet).Duration(args[1].(string), args[2].(time.Duration), args[3].(string))
	p.Ret(4, ret0)
}

func toType1(v interface{}) flag.Value {
	if v == nil {
		return nil
	}
	return v.(flag.Value)
}

func execmFlagSetVar(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	args[0].(*flag.FlagSet).Var(toType1(args[1]), args[2].(string), args[3].(string))
	p.PopN(4)
}

func execmFlagSetParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*flag.FlagSet).Parse(args[1].([]string))
	p.Ret(2, ret0)
}

func execmFlagSetParsed(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*flag.FlagSet).Parsed()
	p.Ret(1, ret0)
}

func execmFlagSetInit(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	args[0].(*flag.FlagSet).Init(args[1].(string), args[2].(flag.ErrorHandling))
	p.PopN(3)
}

func execFloat64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Float64(args[0].(string), args[1].(float64), args[2].(string))
	p.Ret(3, ret0)
}

func execFloat64Var(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.Float64Var(args[0].(*float64), args[1].(string), args[2].(float64), args[3].(string))
	p.PopN(4)
}

func execiGetterGet(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(flag.Getter).Get()
	p.Ret(1, ret0)
}

func execInt(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Int(args[0].(string), args[1].(int), args[2].(string))
	p.Ret(3, ret0)
}

func execInt64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Int64(args[0].(string), args[1].(int64), args[2].(string))
	p.Ret(3, ret0)
}

func execInt64Var(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.Int64Var(args[0].(*int64), args[1].(string), args[2].(int64), args[3].(string))
	p.PopN(4)
}

func execIntVar(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.IntVar(args[0].(*int), args[1].(string), args[2].(int), args[3].(string))
	p.PopN(4)
}

func execLookup(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := flag.Lookup(args[0].(string))
	p.Ret(1, ret0)
}

func execNArg(_ int, p *gop.Context) {
	ret0 := flag.NArg()
	p.Ret(0, ret0)
}

func execNFlag(_ int, p *gop.Context) {
	ret0 := flag.NFlag()
	p.Ret(0, ret0)
}

func execNewFlagSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := flag.NewFlagSet(args[0].(string), args[1].(flag.ErrorHandling))
	p.Ret(2, ret0)
}

func execParse(_ int, p *gop.Context) {
	flag.Parse()
	p.PopN(0)
}

func execParsed(_ int, p *gop.Context) {
	ret0 := flag.Parsed()
	p.Ret(0, ret0)
}

func execPrintDefaults(_ int, p *gop.Context) {
	flag.PrintDefaults()
	p.PopN(0)
}

func execSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := flag.Set(args[0].(string), args[1].(string))
	p.Ret(2, ret0)
}

func execString(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.String(args[0].(string), args[1].(string), args[2].(string))
	p.Ret(3, ret0)
}

func execStringVar(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.StringVar(args[0].(*string), args[1].(string), args[2].(string), args[3].(string))
	p.PopN(4)
}

func execUint(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Uint(args[0].(string), args[1].(uint), args[2].(string))
	p.Ret(3, ret0)
}

func execUint64(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := flag.Uint64(args[0].(string), args[1].(uint64), args[2].(string))
	p.Ret(3, ret0)
}

func execUint64Var(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.Uint64Var(args[0].(*uint64), args[1].(string), args[2].(uint64), args[3].(string))
	p.PopN(4)
}

func execUintVar(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	flag.UintVar(args[0].(*uint), args[1].(string), args[2].(uint), args[3].(string))
	p.PopN(4)
}

func execUnquoteUsage(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := flag.UnquoteUsage(args[0].(*flag.Flag))
	p.Ret(1, ret0, ret1)
}

func execiValueSet(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(flag.Value).Set(args[1].(string))
	p.Ret(2, ret0)
}

func execiValueString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(flag.Value).String()
	p.Ret(1, ret0)
}

func execVar(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	flag.Var(toType1(args[0]), args[1].(string), args[2].(string))
	p.PopN(3)
}

func execVisit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	flag.Visit(args[0].(func(*flag.Flag)))
	p.PopN(1)
}

func execVisitAll(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	flag.VisitAll(args[0].(func(*flag.Flag)))
	p.PopN(1)
}

// I is a Go package instance.
var I = gop.NewGoPackage("flag")

func init() {
	I.RegisterFuncs(
		I.Func("Arg", flag.Arg, execArg),
		I.Func("Args", flag.Args, execArgs),
		I.Func("Bool", flag.Bool, execBool),
		I.Func("BoolVar", flag.BoolVar, execBoolVar),
		I.Func("Duration", flag.Duration, execDuration),
		I.Func("DurationVar", flag.DurationVar, execDurationVar),
		I.Func("(*FlagSet).Output", (*flag.FlagSet).Output, execmFlagSetOutput),
		I.Func("(*FlagSet).Name", (*flag.FlagSet).Name, execmFlagSetName),
		I.Func("(*FlagSet).ErrorHandling", (*flag.FlagSet).ErrorHandling, execmFlagSetErrorHandling),
		I.Func("(*FlagSet).SetOutput", (*flag.FlagSet).SetOutput, execmFlagSetSetOutput),
		I.Func("(*FlagSet).VisitAll", (*flag.FlagSet).VisitAll, execmFlagSetVisitAll),
		I.Func("(*FlagSet).Visit", (*flag.FlagSet).Visit, execmFlagSetVisit),
		I.Func("(*FlagSet).Lookup", (*flag.FlagSet).Lookup, execmFlagSetLookup),
		I.Func("(*FlagSet).Set", (*flag.FlagSet).Set, execmFlagSetSet),
		I.Func("(*FlagSet).PrintDefaults", (*flag.FlagSet).PrintDefaults, execmFlagSetPrintDefaults),
		I.Func("(*FlagSet).NFlag", (*flag.FlagSet).NFlag, execmFlagSetNFlag),
		I.Func("(*FlagSet).Arg", (*flag.FlagSet).Arg, execmFlagSetArg),
		I.Func("(*FlagSet).NArg", (*flag.FlagSet).NArg, execmFlagSetNArg),
		I.Func("(*FlagSet).Args", (*flag.FlagSet).Args, execmFlagSetArgs),
		I.Func("(*FlagSet).BoolVar", (*flag.FlagSet).BoolVar, execmFlagSetBoolVar),
		I.Func("(*FlagSet).Bool", (*flag.FlagSet).Bool, execmFlagSetBool),
		I.Func("(*FlagSet).IntVar", (*flag.FlagSet).IntVar, execmFlagSetIntVar),
		I.Func("(*FlagSet).Int", (*flag.FlagSet).Int, execmFlagSetInt),
		I.Func("(*FlagSet).Int64Var", (*flag.FlagSet).Int64Var, execmFlagSetInt64Var),
		I.Func("(*FlagSet).Int64", (*flag.FlagSet).Int64, execmFlagSetInt64),
		I.Func("(*FlagSet).UintVar", (*flag.FlagSet).UintVar, execmFlagSetUintVar),
		I.Func("(*FlagSet).Uint", (*flag.FlagSet).Uint, execmFlagSetUint),
		I.Func("(*FlagSet).Uint64Var", (*flag.FlagSet).Uint64Var, execmFlagSetUint64Var),
		I.Func("(*FlagSet).Uint64", (*flag.FlagSet).Uint64, execmFlagSetUint64),
		I.Func("(*FlagSet).StringVar", (*flag.FlagSet).StringVar, execmFlagSetStringVar),
		I.Func("(*FlagSet).String", (*flag.FlagSet).String, execmFlagSetString),
		I.Func("(*FlagSet).Float64Var", (*flag.FlagSet).Float64Var, execmFlagSetFloat64Var),
		I.Func("(*FlagSet).Float64", (*flag.FlagSet).Float64, execmFlagSetFloat64),
		I.Func("(*FlagSet).DurationVar", (*flag.FlagSet).DurationVar, execmFlagSetDurationVar),
		I.Func("(*FlagSet).Duration", (*flag.FlagSet).Duration, execmFlagSetDuration),
		I.Func("(*FlagSet).Var", (*flag.FlagSet).Var, execmFlagSetVar),
		I.Func("(*FlagSet).Parse", (*flag.FlagSet).Parse, execmFlagSetParse),
		I.Func("(*FlagSet).Parsed", (*flag.FlagSet).Parsed, execmFlagSetParsed),
		I.Func("(*FlagSet).Init", (*flag.FlagSet).Init, execmFlagSetInit),
		I.Func("Float64", flag.Float64, execFloat64),
		I.Func("Float64Var", flag.Float64Var, execFloat64Var),
		I.Func("(Getter).Get", (flag.Getter).Get, execiGetterGet),
		I.Func("(Getter).Set", (flag.Value).Set, execiValueSet),
		I.Func("(Getter).String", (flag.Value).String, execiValueString),
		I.Func("Int", flag.Int, execInt),
		I.Func("Int64", flag.Int64, execInt64),
		I.Func("Int64Var", flag.Int64Var, execInt64Var),
		I.Func("IntVar", flag.IntVar, execIntVar),
		I.Func("Lookup", flag.Lookup, execLookup),
		I.Func("NArg", flag.NArg, execNArg),
		I.Func("NFlag", flag.NFlag, execNFlag),
		I.Func("NewFlagSet", flag.NewFlagSet, execNewFlagSet),
		I.Func("Parse", flag.Parse, execParse),
		I.Func("Parsed", flag.Parsed, execParsed),
		I.Func("PrintDefaults", flag.PrintDefaults, execPrintDefaults),
		I.Func("Set", flag.Set, execSet),
		I.Func("String", flag.String, execString),
		I.Func("StringVar", flag.StringVar, execStringVar),
		I.Func("Uint", flag.Uint, execUint),
		I.Func("Uint64", flag.Uint64, execUint64),
		I.Func("Uint64Var", flag.Uint64Var, execUint64Var),
		I.Func("UintVar", flag.UintVar, execUintVar),
		I.Func("UnquoteUsage", flag.UnquoteUsage, execUnquoteUsage),
		I.Func("(Value).Set", (flag.Value).Set, execiValueSet),
		I.Func("(Value).String", (flag.Value).String, execiValueString),
		I.Func("Var", flag.Var, execVar),
		I.Func("Visit", flag.Visit, execVisit),
		I.Func("VisitAll", flag.VisitAll, execVisitAll),
	)
	I.RegisterVars(
		I.Var("CommandLine", &flag.CommandLine),
		I.Var("ErrHelp", &flag.ErrHelp),
		I.Var("Usage", &flag.Usage),
	)
	I.RegisterTypes(
		I.Type("ErrorHandling", reflect.TypeOf((*flag.ErrorHandling)(nil)).Elem()),
		I.Type("Flag", reflect.TypeOf((*flag.Flag)(nil)).Elem()),
		I.Type("FlagSet", reflect.TypeOf((*flag.FlagSet)(nil)).Elem()),
		I.Type("Getter", reflect.TypeOf((*flag.Getter)(nil)).Elem()),
		I.Type("Value", reflect.TypeOf((*flag.Value)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ContinueOnError", qspec.Int, flag.ContinueOnError),
		I.Const("ExitOnError", qspec.Int, flag.ExitOnError),
		I.Const("PanicOnError", qspec.Int, flag.PanicOnError),
	)
}

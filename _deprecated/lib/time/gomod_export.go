// Package time provide Go+ "time" package, as "time" package in Go.
package time

import (
	reflect "reflect"
	time "time"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func execAfter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := time.After(args[0].(time.Duration))
	p.Ret(1, ret0)
}

func execAfterFunc(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := time.AfterFunc(args[0].(time.Duration), args[1].(func()))
	p.Ret(2, ret0)
}

func execDate(_ int, p *gop.Context) {
	args := p.GetArgs(8)
	ret0 := time.Date(args[0].(int), args[1].(time.Month), args[2].(int), args[3].(int), args[4].(int), args[5].(int), args[6].(int), args[7].(*time.Location))
	p.Ret(8, ret0)
}

func execmDurationString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).String()
	p.Ret(1, ret0)
}

func execmDurationNanoseconds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).Nanoseconds()
	p.Ret(1, ret0)
}

func execmDurationMicroseconds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).Microseconds()
	p.Ret(1, ret0)
}

func execmDurationMilliseconds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).Milliseconds()
	p.Ret(1, ret0)
}

func execmDurationSeconds(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).Seconds()
	p.Ret(1, ret0)
}

func execmDurationMinutes(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).Minutes()
	p.Ret(1, ret0)
}

func execmDurationHours(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Duration).Hours()
	p.Ret(1, ret0)
}

func execmDurationTruncate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Duration).Truncate(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execmDurationRound(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Duration).Round(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execFixedZone(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := time.FixedZone(args[0].(string), args[1].(int))
	p.Ret(2, ret0)
}

func execLoadLocation(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := time.LoadLocation(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execLoadLocationFromTZData(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := time.LoadLocationFromTZData(args[0].(string), args[1].([]byte))
	p.Ret(2, ret0, ret1)
}

func execmLocationString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*time.Location).String()
	p.Ret(1, ret0)
}

func execmMonthString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Month).String()
	p.Ret(1, ret0)
}

func execNewTicker(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := time.NewTicker(args[0].(time.Duration))
	p.Ret(1, ret0)
}

func execNewTimer(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := time.NewTimer(args[0].(time.Duration))
	p.Ret(1, ret0)
}

func execNow(_ int, p *gop.Context) {
	ret0 := time.Now()
	p.Ret(0, ret0)
}

func execParse(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0, ret1 := time.Parse(args[0].(string), args[1].(string))
	p.Ret(2, ret0, ret1)
}

func execParseDuration(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := time.ParseDuration(args[0].(string))
	p.Ret(1, ret0, ret1)
}

func execmParseErrorError(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*time.ParseError).Error()
	p.Ret(1, ret0)
}

func execParseInLocation(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0, ret1 := time.ParseInLocation(args[0].(string), args[1].(string), args[2].(*time.Location))
	p.Ret(3, ret0, ret1)
}

func execSince(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := time.Since(args[0].(time.Time))
	p.Ret(1, ret0)
}

func execSleep(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	time.Sleep(args[0].(time.Duration))
	p.PopN(1)
}

func execTick(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := time.Tick(args[0].(time.Duration))
	p.Ret(1, ret0)
}

func execmTickerStop(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	args[0].(*time.Ticker).Stop()
	p.PopN(1)
}

func execmTimeString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).String()
	p.Ret(1, ret0)
}

func execmTimeFormat(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Format(args[1].(string))
	p.Ret(2, ret0)
}

func execmTimeAppendFormat(_ int, p *gop.Context) {
	args := p.GetArgs(3)
	ret0 := args[0].(time.Time).AppendFormat(args[1].([]byte), args[2].(string))
	p.Ret(3, ret0)
}

func execmTimeAfter(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).After(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmTimeBefore(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Before(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmTimeEqual(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Equal(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmTimeIsZero(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).IsZero()
	p.Ret(1, ret0)
}

func execmTimeDate(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(time.Time).Date()
	p.Ret(1, ret0, ret1, ret2)
}

func execmTimeYear(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Year()
	p.Ret(1, ret0)
}

func execmTimeMonth(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Month()
	p.Ret(1, ret0)
}

func execmTimeDay(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Day()
	p.Ret(1, ret0)
}

func execmTimeWeekday(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Weekday()
	p.Ret(1, ret0)
}

func execmTimeISOWeek(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(time.Time).ISOWeek()
	p.Ret(1, ret0, ret1)
}

func execmTimeClock(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1, ret2 := args[0].(time.Time).Clock()
	p.Ret(1, ret0, ret1, ret2)
}

func execmTimeHour(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Hour()
	p.Ret(1, ret0)
}

func execmTimeMinute(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Minute()
	p.Ret(1, ret0)
}

func execmTimeSecond(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Second()
	p.Ret(1, ret0)
}

func execmTimeNanosecond(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Nanosecond()
	p.Ret(1, ret0)
}

func execmTimeYearDay(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).YearDay()
	p.Ret(1, ret0)
}

func execmTimeAdd(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Add(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execmTimeSub(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Sub(args[1].(time.Time))
	p.Ret(2, ret0)
}

func execmTimeAddDate(_ int, p *gop.Context) {
	args := p.GetArgs(4)
	ret0 := args[0].(time.Time).AddDate(args[1].(int), args[2].(int), args[3].(int))
	p.Ret(4, ret0)
}

func execmTimeUTC(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).UTC()
	p.Ret(1, ret0)
}

func execmTimeLocal(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Local()
	p.Ret(1, ret0)
}

func execmTimeIn(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).In(args[1].(*time.Location))
	p.Ret(2, ret0)
}

func execmTimeLocation(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Location()
	p.Ret(1, ret0)
}

func execmTimeZone(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(time.Time).Zone()
	p.Ret(1, ret0, ret1)
}

func execmTimeUnix(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).Unix()
	p.Ret(1, ret0)
}

func execmTimeUnixNano(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Time).UnixNano()
	p.Ret(1, ret0)
}

func execmTimeMarshalBinary(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(time.Time).MarshalBinary()
	p.Ret(1, ret0, ret1)
}

func execmTimeUnmarshalBinary(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*time.Time).UnmarshalBinary(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmTimeGobEncode(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(time.Time).GobEncode()
	p.Ret(1, ret0, ret1)
}

func execmTimeGobDecode(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*time.Time).GobDecode(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmTimeMarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(time.Time).MarshalJSON()
	p.Ret(1, ret0, ret1)
}

func execmTimeUnmarshalJSON(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*time.Time).UnmarshalJSON(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmTimeMarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0, ret1 := args[0].(time.Time).MarshalText()
	p.Ret(1, ret0, ret1)
}

func execmTimeUnmarshalText(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*time.Time).UnmarshalText(args[1].([]byte))
	p.Ret(2, ret0)
}

func execmTimeTruncate(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Truncate(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execmTimeRound(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(time.Time).Round(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execmTimerStop(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(*time.Timer).Stop()
	p.Ret(1, ret0)
}

func execmTimerReset(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(*time.Timer).Reset(args[1].(time.Duration))
	p.Ret(2, ret0)
}

func execUnix(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := time.Unix(args[0].(int64), args[1].(int64))
	p.Ret(2, ret0)
}

func execUntil(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := time.Until(args[0].(time.Time))
	p.Ret(1, ret0)
}

func execmWeekdayString(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := args[0].(time.Weekday).String()
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("time")

func init() {
	I.RegisterFuncs(
		I.Func("After", time.After, execAfter),
		I.Func("AfterFunc", time.AfterFunc, execAfterFunc),
		I.Func("Date", time.Date, execDate),
		I.Func("(Duration).String", (time.Duration).String, execmDurationString),
		I.Func("(Duration).Nanoseconds", (time.Duration).Nanoseconds, execmDurationNanoseconds),
		I.Func("(Duration).Microseconds", (time.Duration).Microseconds, execmDurationMicroseconds),
		I.Func("(Duration).Milliseconds", (time.Duration).Milliseconds, execmDurationMilliseconds),
		I.Func("(Duration).Seconds", (time.Duration).Seconds, execmDurationSeconds),
		I.Func("(Duration).Minutes", (time.Duration).Minutes, execmDurationMinutes),
		I.Func("(Duration).Hours", (time.Duration).Hours, execmDurationHours),
		I.Func("(Duration).Truncate", (time.Duration).Truncate, execmDurationTruncate),
		I.Func("(Duration).Round", (time.Duration).Round, execmDurationRound),
		I.Func("FixedZone", time.FixedZone, execFixedZone),
		I.Func("LoadLocation", time.LoadLocation, execLoadLocation),
		I.Func("LoadLocationFromTZData", time.LoadLocationFromTZData, execLoadLocationFromTZData),
		I.Func("(*Location).String", (*time.Location).String, execmLocationString),
		I.Func("(Month).String", (time.Month).String, execmMonthString),
		I.Func("NewTicker", time.NewTicker, execNewTicker),
		I.Func("NewTimer", time.NewTimer, execNewTimer),
		I.Func("Now", time.Now, execNow),
		I.Func("Parse", time.Parse, execParse),
		I.Func("ParseDuration", time.ParseDuration, execParseDuration),
		I.Func("(*ParseError).Error", (*time.ParseError).Error, execmParseErrorError),
		I.Func("ParseInLocation", time.ParseInLocation, execParseInLocation),
		I.Func("Since", time.Since, execSince),
		I.Func("Sleep", time.Sleep, execSleep),
		I.Func("Tick", time.Tick, execTick),
		I.Func("(*Ticker).Stop", (*time.Ticker).Stop, execmTickerStop),
		I.Func("(Time).String", (time.Time).String, execmTimeString),
		I.Func("(Time).Format", (time.Time).Format, execmTimeFormat),
		I.Func("(Time).AppendFormat", (time.Time).AppendFormat, execmTimeAppendFormat),
		I.Func("(Time).After", (time.Time).After, execmTimeAfter),
		I.Func("(Time).Before", (time.Time).Before, execmTimeBefore),
		I.Func("(Time).Equal", (time.Time).Equal, execmTimeEqual),
		I.Func("(Time).IsZero", (time.Time).IsZero, execmTimeIsZero),
		I.Func("(Time).Date", (time.Time).Date, execmTimeDate),
		I.Func("(Time).Year", (time.Time).Year, execmTimeYear),
		I.Func("(Time).Month", (time.Time).Month, execmTimeMonth),
		I.Func("(Time).Day", (time.Time).Day, execmTimeDay),
		I.Func("(Time).Weekday", (time.Time).Weekday, execmTimeWeekday),
		I.Func("(Time).ISOWeek", (time.Time).ISOWeek, execmTimeISOWeek),
		I.Func("(Time).Clock", (time.Time).Clock, execmTimeClock),
		I.Func("(Time).Hour", (time.Time).Hour, execmTimeHour),
		I.Func("(Time).Minute", (time.Time).Minute, execmTimeMinute),
		I.Func("(Time).Second", (time.Time).Second, execmTimeSecond),
		I.Func("(Time).Nanosecond", (time.Time).Nanosecond, execmTimeNanosecond),
		I.Func("(Time).YearDay", (time.Time).YearDay, execmTimeYearDay),
		I.Func("(Time).Add", (time.Time).Add, execmTimeAdd),
		I.Func("(Time).Sub", (time.Time).Sub, execmTimeSub),
		I.Func("(Time).AddDate", (time.Time).AddDate, execmTimeAddDate),
		I.Func("(Time).UTC", (time.Time).UTC, execmTimeUTC),
		I.Func("(Time).Local", (time.Time).Local, execmTimeLocal),
		I.Func("(Time).In", (time.Time).In, execmTimeIn),
		I.Func("(Time).Location", (time.Time).Location, execmTimeLocation),
		I.Func("(Time).Zone", (time.Time).Zone, execmTimeZone),
		I.Func("(Time).Unix", (time.Time).Unix, execmTimeUnix),
		I.Func("(Time).UnixNano", (time.Time).UnixNano, execmTimeUnixNano),
		I.Func("(Time).MarshalBinary", (time.Time).MarshalBinary, execmTimeMarshalBinary),
		I.Func("(*Time).UnmarshalBinary", (*time.Time).UnmarshalBinary, execmTimeUnmarshalBinary),
		I.Func("(Time).GobEncode", (time.Time).GobEncode, execmTimeGobEncode),
		I.Func("(*Time).GobDecode", (*time.Time).GobDecode, execmTimeGobDecode),
		I.Func("(Time).MarshalJSON", (time.Time).MarshalJSON, execmTimeMarshalJSON),
		I.Func("(*Time).UnmarshalJSON", (*time.Time).UnmarshalJSON, execmTimeUnmarshalJSON),
		I.Func("(Time).MarshalText", (time.Time).MarshalText, execmTimeMarshalText),
		I.Func("(*Time).UnmarshalText", (*time.Time).UnmarshalText, execmTimeUnmarshalText),
		I.Func("(Time).Truncate", (time.Time).Truncate, execmTimeTruncate),
		I.Func("(Time).Round", (time.Time).Round, execmTimeRound),
		I.Func("(*Timer).Stop", (*time.Timer).Stop, execmTimerStop),
		I.Func("(*Timer).Reset", (*time.Timer).Reset, execmTimerReset),
		I.Func("Unix", time.Unix, execUnix),
		I.Func("Until", time.Until, execUntil),
		I.Func("(Weekday).String", (time.Weekday).String, execmWeekdayString),
	)
	I.RegisterVars(
		I.Var("Local", &time.Local),
		I.Var("UTC", &time.UTC),
	)
	I.RegisterTypes(
		I.Type("Duration", reflect.TypeOf((*time.Duration)(nil)).Elem()),
		I.Type("Location", reflect.TypeOf((*time.Location)(nil)).Elem()),
		I.Type("Month", reflect.TypeOf((*time.Month)(nil)).Elem()),
		I.Type("ParseError", reflect.TypeOf((*time.ParseError)(nil)).Elem()),
		I.Type("Ticker", reflect.TypeOf((*time.Ticker)(nil)).Elem()),
		I.Type("Time", reflect.TypeOf((*time.Time)(nil)).Elem()),
		I.Type("Timer", reflect.TypeOf((*time.Timer)(nil)).Elem()),
		I.Type("Weekday", reflect.TypeOf((*time.Weekday)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("ANSIC", qspec.ConstBoundString, time.ANSIC),
		I.Const("April", qspec.Int, time.April),
		I.Const("August", qspec.Int, time.August),
		I.Const("December", qspec.Int, time.December),
		I.Const("February", qspec.Int, time.February),
		I.Const("Friday", qspec.Int, time.Friday),
		I.Const("Hour", qspec.Int64, time.Hour),
		I.Const("January", qspec.Int, time.January),
		I.Const("July", qspec.Int, time.July),
		I.Const("June", qspec.Int, time.June),
		I.Const("Kitchen", qspec.ConstBoundString, time.Kitchen),
		I.Const("March", qspec.Int, time.March),
		I.Const("May", qspec.Int, time.May),
		I.Const("Microsecond", qspec.Int64, time.Microsecond),
		I.Const("Millisecond", qspec.Int64, time.Millisecond),
		I.Const("Minute", qspec.Int64, time.Minute),
		I.Const("Monday", qspec.Int, time.Monday),
		I.Const("Nanosecond", qspec.Int64, time.Nanosecond),
		I.Const("November", qspec.Int, time.November),
		I.Const("October", qspec.Int, time.October),
		I.Const("RFC1123", qspec.ConstBoundString, time.RFC1123),
		I.Const("RFC1123Z", qspec.ConstBoundString, time.RFC1123Z),
		I.Const("RFC3339", qspec.ConstBoundString, time.RFC3339),
		I.Const("RFC3339Nano", qspec.ConstBoundString, time.RFC3339Nano),
		I.Const("RFC822", qspec.ConstBoundString, time.RFC822),
		I.Const("RFC822Z", qspec.ConstBoundString, time.RFC822Z),
		I.Const("RFC850", qspec.ConstBoundString, time.RFC850),
		I.Const("RubyDate", qspec.ConstBoundString, time.RubyDate),
		I.Const("Saturday", qspec.Int, time.Saturday),
		I.Const("Second", qspec.Int64, time.Second),
		I.Const("September", qspec.Int, time.September),
		I.Const("Stamp", qspec.ConstBoundString, time.Stamp),
		I.Const("StampMicro", qspec.ConstBoundString, time.StampMicro),
		I.Const("StampMilli", qspec.ConstBoundString, time.StampMilli),
		I.Const("StampNano", qspec.ConstBoundString, time.StampNano),
		I.Const("Sunday", qspec.Int, time.Sunday),
		I.Const("Thursday", qspec.Int, time.Thursday),
		I.Const("Tuesday", qspec.Int, time.Tuesday),
		I.Const("UnixDate", qspec.ConstBoundString, time.UnixDate),
		I.Const("Wednesday", qspec.Int, time.Wednesday),
	)
}

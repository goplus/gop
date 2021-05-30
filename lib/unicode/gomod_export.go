// Package unicode provide Go+ "unicode" package, as "unicode" package in Go.
package unicode

import (
	reflect "reflect"
	unicode "unicode"

	gop "github.com/goplus/gop"
	qspec "github.com/goplus/gop/exec.spec"
)

func toSlice0(args []interface{}) []*unicode.RangeTable {
	ret := make([]*unicode.RangeTable, len(args))
	for i, arg := range args {
		ret[i] = arg.(*unicode.RangeTable)
	}
	return ret
}

func execIn(arity int, p *gop.Context) {
	args := p.GetArgs(arity)
	ret0 := unicode.In(args[0].(rune), toSlice0(args[1:])...)
	p.Ret(arity, ret0)
}

func execIs(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := unicode.Is(args[0].(*unicode.RangeTable), args[1].(rune))
	p.Ret(2, ret0)
}

func execIsControl(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsControl(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsDigit(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsDigit(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsGraphic(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsGraphic(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsLetter(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsLetter(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsLower(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsLower(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsMark(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsMark(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsNumber(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsNumber(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsOneOf(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := unicode.IsOneOf(args[0].([]*unicode.RangeTable), args[1].(rune))
	p.Ret(2, ret0)
}

func execIsPrint(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsPrint(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsPunct(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsPunct(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsSpace(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsSpace(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsSymbol(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsSymbol(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsTitle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsTitle(args[0].(rune))
	p.Ret(1, ret0)
}

func execIsUpper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.IsUpper(args[0].(rune))
	p.Ret(1, ret0)
}

func execSimpleFold(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.SimpleFold(args[0].(rune))
	p.Ret(1, ret0)
}

func execmSpecialCaseToUpper(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(unicode.SpecialCase).ToUpper(args[1].(rune))
	p.Ret(2, ret0)
}

func execmSpecialCaseToTitle(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(unicode.SpecialCase).ToTitle(args[1].(rune))
	p.Ret(2, ret0)
}

func execmSpecialCaseToLower(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := args[0].(unicode.SpecialCase).ToLower(args[1].(rune))
	p.Ret(2, ret0)
}

func execTo(_ int, p *gop.Context) {
	args := p.GetArgs(2)
	ret0 := unicode.To(args[0].(int), args[1].(rune))
	p.Ret(2, ret0)
}

func execToLower(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.ToLower(args[0].(rune))
	p.Ret(1, ret0)
}

func execToTitle(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.ToTitle(args[0].(rune))
	p.Ret(1, ret0)
}

func execToUpper(_ int, p *gop.Context) {
	args := p.GetArgs(1)
	ret0 := unicode.ToUpper(args[0].(rune))
	p.Ret(1, ret0)
}

// I is a Go package instance.
var I = gop.NewGoPackage("unicode")

func init() {
	I.RegisterFuncs(
		I.Func("Is", unicode.Is, execIs),
		I.Func("IsControl", unicode.IsControl, execIsControl),
		I.Func("IsDigit", unicode.IsDigit, execIsDigit),
		I.Func("IsGraphic", unicode.IsGraphic, execIsGraphic),
		I.Func("IsLetter", unicode.IsLetter, execIsLetter),
		I.Func("IsLower", unicode.IsLower, execIsLower),
		I.Func("IsMark", unicode.IsMark, execIsMark),
		I.Func("IsNumber", unicode.IsNumber, execIsNumber),
		I.Func("IsOneOf", unicode.IsOneOf, execIsOneOf),
		I.Func("IsPrint", unicode.IsPrint, execIsPrint),
		I.Func("IsPunct", unicode.IsPunct, execIsPunct),
		I.Func("IsSpace", unicode.IsSpace, execIsSpace),
		I.Func("IsSymbol", unicode.IsSymbol, execIsSymbol),
		I.Func("IsTitle", unicode.IsTitle, execIsTitle),
		I.Func("IsUpper", unicode.IsUpper, execIsUpper),
		I.Func("SimpleFold", unicode.SimpleFold, execSimpleFold),
		I.Func("(SpecialCase).ToUpper", (unicode.SpecialCase).ToUpper, execmSpecialCaseToUpper),
		I.Func("(SpecialCase).ToTitle", (unicode.SpecialCase).ToTitle, execmSpecialCaseToTitle),
		I.Func("(SpecialCase).ToLower", (unicode.SpecialCase).ToLower, execmSpecialCaseToLower),
		I.Func("To", unicode.To, execTo),
		I.Func("ToLower", unicode.ToLower, execToLower),
		I.Func("ToTitle", unicode.ToTitle, execToTitle),
		I.Func("ToUpper", unicode.ToUpper, execToUpper),
	)
	I.RegisterFuncvs(
		I.Funcv("In", unicode.In, execIn),
	)
	I.RegisterVars(
		I.Var("ASCII_Hex_Digit", &unicode.ASCII_Hex_Digit),
		I.Var("Adlam", &unicode.Adlam),
		I.Var("Ahom", &unicode.Ahom),
		I.Var("Anatolian_Hieroglyphs", &unicode.Anatolian_Hieroglyphs),
		I.Var("Arabic", &unicode.Arabic),
		I.Var("Armenian", &unicode.Armenian),
		I.Var("Avestan", &unicode.Avestan),
		I.Var("AzeriCase", &unicode.AzeriCase),
		I.Var("Balinese", &unicode.Balinese),
		I.Var("Bamum", &unicode.Bamum),
		I.Var("Bassa_Vah", &unicode.Bassa_Vah),
		I.Var("Batak", &unicode.Batak),
		I.Var("Bengali", &unicode.Bengali),
		I.Var("Bhaiksuki", &unicode.Bhaiksuki),
		I.Var("Bidi_Control", &unicode.Bidi_Control),
		I.Var("Bopomofo", &unicode.Bopomofo),
		I.Var("Brahmi", &unicode.Brahmi),
		I.Var("Braille", &unicode.Braille),
		I.Var("Buginese", &unicode.Buginese),
		I.Var("Buhid", &unicode.Buhid),
		I.Var("C", &unicode.C),
		I.Var("Canadian_Aboriginal", &unicode.Canadian_Aboriginal),
		I.Var("Carian", &unicode.Carian),
		I.Var("CaseRanges", &unicode.CaseRanges),
		I.Var("Categories", &unicode.Categories),
		I.Var("Caucasian_Albanian", &unicode.Caucasian_Albanian),
		I.Var("Cc", &unicode.Cc),
		I.Var("Cf", &unicode.Cf),
		I.Var("Chakma", &unicode.Chakma),
		I.Var("Cham", &unicode.Cham),
		I.Var("Cherokee", &unicode.Cherokee),
		I.Var("Co", &unicode.Co),
		I.Var("Common", &unicode.Common),
		I.Var("Coptic", &unicode.Coptic),
		I.Var("Cs", &unicode.Cs),
		I.Var("Cuneiform", &unicode.Cuneiform),
		I.Var("Cypriot", &unicode.Cypriot),
		I.Var("Cyrillic", &unicode.Cyrillic),
		I.Var("Dash", &unicode.Dash),
		I.Var("Deprecated", &unicode.Deprecated),
		I.Var("Deseret", &unicode.Deseret),
		I.Var("Devanagari", &unicode.Devanagari),
		I.Var("Diacritic", &unicode.Diacritic),
		I.Var("Digit", &unicode.Digit),
		I.Var("Dogra", &unicode.Dogra),
		I.Var("Duployan", &unicode.Duployan),
		I.Var("Egyptian_Hieroglyphs", &unicode.Egyptian_Hieroglyphs),
		I.Var("Elbasan", &unicode.Elbasan),
		I.Var("Elymaic", &unicode.Elymaic),
		I.Var("Ethiopic", &unicode.Ethiopic),
		I.Var("Extender", &unicode.Extender),
		I.Var("FoldCategory", &unicode.FoldCategory),
		I.Var("FoldScript", &unicode.FoldScript),
		I.Var("Georgian", &unicode.Georgian),
		I.Var("Glagolitic", &unicode.Glagolitic),
		I.Var("Gothic", &unicode.Gothic),
		I.Var("Grantha", &unicode.Grantha),
		I.Var("GraphicRanges", &unicode.GraphicRanges),
		I.Var("Greek", &unicode.Greek),
		I.Var("Gujarati", &unicode.Gujarati),
		I.Var("Gunjala_Gondi", &unicode.Gunjala_Gondi),
		I.Var("Gurmukhi", &unicode.Gurmukhi),
		I.Var("Han", &unicode.Han),
		I.Var("Hangul", &unicode.Hangul),
		I.Var("Hanifi_Rohingya", &unicode.Hanifi_Rohingya),
		I.Var("Hanunoo", &unicode.Hanunoo),
		I.Var("Hatran", &unicode.Hatran),
		I.Var("Hebrew", &unicode.Hebrew),
		I.Var("Hex_Digit", &unicode.Hex_Digit),
		I.Var("Hiragana", &unicode.Hiragana),
		I.Var("Hyphen", &unicode.Hyphen),
		I.Var("IDS_Binary_Operator", &unicode.IDS_Binary_Operator),
		I.Var("IDS_Trinary_Operator", &unicode.IDS_Trinary_Operator),
		I.Var("Ideographic", &unicode.Ideographic),
		I.Var("Imperial_Aramaic", &unicode.Imperial_Aramaic),
		I.Var("Inherited", &unicode.Inherited),
		I.Var("Inscriptional_Pahlavi", &unicode.Inscriptional_Pahlavi),
		I.Var("Inscriptional_Parthian", &unicode.Inscriptional_Parthian),
		I.Var("Javanese", &unicode.Javanese),
		I.Var("Join_Control", &unicode.Join_Control),
		I.Var("Kaithi", &unicode.Kaithi),
		I.Var("Kannada", &unicode.Kannada),
		I.Var("Katakana", &unicode.Katakana),
		I.Var("Kayah_Li", &unicode.Kayah_Li),
		I.Var("Kharoshthi", &unicode.Kharoshthi),
		I.Var("Khmer", &unicode.Khmer),
		I.Var("Khojki", &unicode.Khojki),
		I.Var("Khudawadi", &unicode.Khudawadi),
		I.Var("L", &unicode.L),
		I.Var("Lao", &unicode.Lao),
		I.Var("Latin", &unicode.Latin),
		I.Var("Lepcha", &unicode.Lepcha),
		I.Var("Letter", &unicode.Letter),
		I.Var("Limbu", &unicode.Limbu),
		I.Var("Linear_A", &unicode.Linear_A),
		I.Var("Linear_B", &unicode.Linear_B),
		I.Var("Lisu", &unicode.Lisu),
		I.Var("Ll", &unicode.Ll),
		I.Var("Lm", &unicode.Lm),
		I.Var("Lo", &unicode.Lo),
		I.Var("Logical_Order_Exception", &unicode.Logical_Order_Exception),
		I.Var("Lower", &unicode.Lower),
		I.Var("Lt", &unicode.Lt),
		I.Var("Lu", &unicode.Lu),
		I.Var("Lycian", &unicode.Lycian),
		I.Var("Lydian", &unicode.Lydian),
		I.Var("M", &unicode.M),
		I.Var("Mahajani", &unicode.Mahajani),
		I.Var("Makasar", &unicode.Makasar),
		I.Var("Malayalam", &unicode.Malayalam),
		I.Var("Mandaic", &unicode.Mandaic),
		I.Var("Manichaean", &unicode.Manichaean),
		I.Var("Marchen", &unicode.Marchen),
		I.Var("Mark", &unicode.Mark),
		I.Var("Masaram_Gondi", &unicode.Masaram_Gondi),
		I.Var("Mc", &unicode.Mc),
		I.Var("Me", &unicode.Me),
		I.Var("Medefaidrin", &unicode.Medefaidrin),
		I.Var("Meetei_Mayek", &unicode.Meetei_Mayek),
		I.Var("Mende_Kikakui", &unicode.Mende_Kikakui),
		I.Var("Meroitic_Cursive", &unicode.Meroitic_Cursive),
		I.Var("Meroitic_Hieroglyphs", &unicode.Meroitic_Hieroglyphs),
		I.Var("Miao", &unicode.Miao),
		I.Var("Mn", &unicode.Mn),
		I.Var("Modi", &unicode.Modi),
		I.Var("Mongolian", &unicode.Mongolian),
		I.Var("Mro", &unicode.Mro),
		I.Var("Multani", &unicode.Multani),
		I.Var("Myanmar", &unicode.Myanmar),
		I.Var("N", &unicode.N),
		I.Var("Nabataean", &unicode.Nabataean),
		I.Var("Nandinagari", &unicode.Nandinagari),
		I.Var("Nd", &unicode.Nd),
		I.Var("New_Tai_Lue", &unicode.New_Tai_Lue),
		I.Var("Newa", &unicode.Newa),
		I.Var("Nko", &unicode.Nko),
		I.Var("Nl", &unicode.Nl),
		I.Var("No", &unicode.No),
		I.Var("Noncharacter_Code_Point", &unicode.Noncharacter_Code_Point),
		I.Var("Number", &unicode.Number),
		I.Var("Nushu", &unicode.Nushu),
		I.Var("Nyiakeng_Puachue_Hmong", &unicode.Nyiakeng_Puachue_Hmong),
		I.Var("Ogham", &unicode.Ogham),
		I.Var("Ol_Chiki", &unicode.Ol_Chiki),
		I.Var("Old_Hungarian", &unicode.Old_Hungarian),
		I.Var("Old_Italic", &unicode.Old_Italic),
		I.Var("Old_North_Arabian", &unicode.Old_North_Arabian),
		I.Var("Old_Permic", &unicode.Old_Permic),
		I.Var("Old_Persian", &unicode.Old_Persian),
		I.Var("Old_Sogdian", &unicode.Old_Sogdian),
		I.Var("Old_South_Arabian", &unicode.Old_South_Arabian),
		I.Var("Old_Turkic", &unicode.Old_Turkic),
		I.Var("Oriya", &unicode.Oriya),
		I.Var("Osage", &unicode.Osage),
		I.Var("Osmanya", &unicode.Osmanya),
		I.Var("Other", &unicode.Other),
		I.Var("Other_Alphabetic", &unicode.Other_Alphabetic),
		I.Var("Other_Default_Ignorable_Code_Point", &unicode.Other_Default_Ignorable_Code_Point),
		I.Var("Other_Grapheme_Extend", &unicode.Other_Grapheme_Extend),
		I.Var("Other_ID_Continue", &unicode.Other_ID_Continue),
		I.Var("Other_ID_Start", &unicode.Other_ID_Start),
		I.Var("Other_Lowercase", &unicode.Other_Lowercase),
		I.Var("Other_Math", &unicode.Other_Math),
		I.Var("Other_Uppercase", &unicode.Other_Uppercase),
		I.Var("P", &unicode.P),
		I.Var("Pahawh_Hmong", &unicode.Pahawh_Hmong),
		I.Var("Palmyrene", &unicode.Palmyrene),
		I.Var("Pattern_Syntax", &unicode.Pattern_Syntax),
		I.Var("Pattern_White_Space", &unicode.Pattern_White_Space),
		I.Var("Pau_Cin_Hau", &unicode.Pau_Cin_Hau),
		I.Var("Pc", &unicode.Pc),
		I.Var("Pd", &unicode.Pd),
		I.Var("Pe", &unicode.Pe),
		I.Var("Pf", &unicode.Pf),
		I.Var("Phags_Pa", &unicode.Phags_Pa),
		I.Var("Phoenician", &unicode.Phoenician),
		I.Var("Pi", &unicode.Pi),
		I.Var("Po", &unicode.Po),
		I.Var("Prepended_Concatenation_Mark", &unicode.Prepended_Concatenation_Mark),
		I.Var("PrintRanges", &unicode.PrintRanges),
		I.Var("Properties", &unicode.Properties),
		I.Var("Ps", &unicode.Ps),
		I.Var("Psalter_Pahlavi", &unicode.Psalter_Pahlavi),
		I.Var("Punct", &unicode.Punct),
		I.Var("Quotation_Mark", &unicode.Quotation_Mark),
		I.Var("Radical", &unicode.Radical),
		I.Var("Regional_Indicator", &unicode.Regional_Indicator),
		I.Var("Rejang", &unicode.Rejang),
		I.Var("Runic", &unicode.Runic),
		I.Var("S", &unicode.S),
		I.Var("STerm", &unicode.STerm),
		I.Var("Samaritan", &unicode.Samaritan),
		I.Var("Saurashtra", &unicode.Saurashtra),
		I.Var("Sc", &unicode.Sc),
		I.Var("Scripts", &unicode.Scripts),
		I.Var("Sentence_Terminal", &unicode.Sentence_Terminal),
		I.Var("Sharada", &unicode.Sharada),
		I.Var("Shavian", &unicode.Shavian),
		I.Var("Siddham", &unicode.Siddham),
		I.Var("SignWriting", &unicode.SignWriting),
		I.Var("Sinhala", &unicode.Sinhala),
		I.Var("Sk", &unicode.Sk),
		I.Var("Sm", &unicode.Sm),
		I.Var("So", &unicode.So),
		I.Var("Soft_Dotted", &unicode.Soft_Dotted),
		I.Var("Sogdian", &unicode.Sogdian),
		I.Var("Sora_Sompeng", &unicode.Sora_Sompeng),
		I.Var("Soyombo", &unicode.Soyombo),
		I.Var("Space", &unicode.Space),
		I.Var("Sundanese", &unicode.Sundanese),
		I.Var("Syloti_Nagri", &unicode.Syloti_Nagri),
		I.Var("Symbol", &unicode.Symbol),
		I.Var("Syriac", &unicode.Syriac),
		I.Var("Tagalog", &unicode.Tagalog),
		I.Var("Tagbanwa", &unicode.Tagbanwa),
		I.Var("Tai_Le", &unicode.Tai_Le),
		I.Var("Tai_Tham", &unicode.Tai_Tham),
		I.Var("Tai_Viet", &unicode.Tai_Viet),
		I.Var("Takri", &unicode.Takri),
		I.Var("Tamil", &unicode.Tamil),
		I.Var("Tangut", &unicode.Tangut),
		I.Var("Telugu", &unicode.Telugu),
		I.Var("Terminal_Punctuation", &unicode.Terminal_Punctuation),
		I.Var("Thaana", &unicode.Thaana),
		I.Var("Thai", &unicode.Thai),
		I.Var("Tibetan", &unicode.Tibetan),
		I.Var("Tifinagh", &unicode.Tifinagh),
		I.Var("Tirhuta", &unicode.Tirhuta),
		I.Var("Title", &unicode.Title),
		I.Var("TurkishCase", &unicode.TurkishCase),
		I.Var("Ugaritic", &unicode.Ugaritic),
		I.Var("Unified_Ideograph", &unicode.Unified_Ideograph),
		I.Var("Upper", &unicode.Upper),
		I.Var("Vai", &unicode.Vai),
		I.Var("Variation_Selector", &unicode.Variation_Selector),
		I.Var("Wancho", &unicode.Wancho),
		I.Var("Warang_Citi", &unicode.Warang_Citi),
		I.Var("White_Space", &unicode.White_Space),
		I.Var("Yi", &unicode.Yi),
		I.Var("Z", &unicode.Z),
		I.Var("Zanabazar_Square", &unicode.Zanabazar_Square),
		I.Var("Zl", &unicode.Zl),
		I.Var("Zp", &unicode.Zp),
		I.Var("Zs", &unicode.Zs),
	)
	I.RegisterTypes(
		I.Type("CaseRange", reflect.TypeOf((*unicode.CaseRange)(nil)).Elem()),
		I.Type("Range16", reflect.TypeOf((*unicode.Range16)(nil)).Elem()),
		I.Type("Range32", reflect.TypeOf((*unicode.Range32)(nil)).Elem()),
		I.Type("RangeTable", reflect.TypeOf((*unicode.RangeTable)(nil)).Elem()),
		I.Type("SpecialCase", reflect.TypeOf((*unicode.SpecialCase)(nil)).Elem()),
	)
	I.RegisterConsts(
		I.Const("LowerCase", qspec.ConstUnboundInt, unicode.LowerCase),
		I.Const("MaxASCII", qspec.ConstBoundRune, unicode.MaxASCII),
		I.Const("MaxCase", qspec.ConstUnboundInt, unicode.MaxCase),
		I.Const("MaxLatin1", qspec.ConstBoundRune, unicode.MaxLatin1),
		I.Const("MaxRune", qspec.ConstBoundRune, unicode.MaxRune),
		I.Const("ReplacementChar", qspec.ConstBoundRune, unicode.ReplacementChar),
		I.Const("TitleCase", qspec.ConstUnboundInt, unicode.TitleCase),
		I.Const("UpperCase", qspec.ConstUnboundInt, unicode.UpperCase),
		I.Const("UpperLower", qspec.ConstBoundRune, unicode.UpperLower),
		I.Const("Version", qspec.ConstBoundString, unicode.Version),
	)
}

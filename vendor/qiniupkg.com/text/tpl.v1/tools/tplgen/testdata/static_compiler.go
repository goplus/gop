package main

import (
	"qiniupkg.com/text/tpl.v1"
)

type StaticCompiler struct {
	Grammar  string
	Marker   func(g tpl.Grammar, mark string) tpl.Grammar
	Init     func()
	Scanner  tpl.Tokener
	ScanMode tpl.ScanMode
}

func (cmplr *StaticCompiler) Cl() tpl.CompileRet {
	var (
		factor = new(tpl.GrVar)
		term   = new(tpl.GrNamed)
		doc    = new(tpl.GrNamed)
	)
	// GrVar factor
	factor.Assign(tpl.Or(cmplr.Marker(tpl.Gr(tpl.FLOAT), "push"), cmplr.Marker(tpl.And(tpl.Gr('-'), factor), "neg"), tpl.And(tpl.Gr('('), doc, tpl.Gr(')')), cmplr.Marker(tpl.And(tpl.Gr(tpl.IDENT), tpl.Gr('('), cmplr.Marker(tpl.List(doc, tpl.Gr(',')), "arity"), tpl.Gr(')')), "call"), tpl.And(tpl.Gr('+'), factor)))
	// GrNamed term
	term.Assign("term", tpl.And(factor, tpl.Repeat0(tpl.Or(cmplr.Marker(tpl.And(tpl.Gr('*'), factor), "mul"), cmplr.Marker(tpl.And(tpl.Gr('/'), factor), "div")))))
	// GrNamed doc
	doc.Assign("doc", tpl.And(term, tpl.Repeat0(tpl.Or(cmplr.Marker(tpl.And(tpl.Gr('+'), term), "add"), cmplr.Marker(tpl.And(tpl.Gr('-'), term), "sub")))))
	ret := tpl.CompileRet{}
	ret.Matcher = &tpl.Matcher{
		Grammar: doc,
		Scanner: cmplr.Scanner,
		Init:    cmplr.Init,
	}
	return ret
}

package extractor

import (
	"errors"
	"strconv"
	"strings"

	"qiniupkg.com/text/tpl.v1"
)

/* -----------------------------------------------------------------------------

leafMark = xxx_Int | xxx_Float | xxx_String | xxx_Text | xxx_Char | xxx_Bool
leafArrayMark = xxx_Ints | xxx_Floats | ... | xxx_Bools
leafArrayMark = xxx_IntArray | xxx_FloatArray | ... | xxx_BoolArray
nodeMark = xxx_Node
nodeArrayMark = xxx_NodeArray

// ---------------------------------------------------------------------------*/

func GetValue(typ string, lit string) interface{} {

	switch typ {
	case "String":
		v, err := strconv.Unquote(lit)
		if err != nil {
			panic(err)
		}
		return v
	case "Int":
		v, err := strconv.ParseInt(lit, 10, 64)
		if err != nil {
			panic(err)
		}
		return v
	case "Float":
		v, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			panic(err)
		}
		return v
	case "Char":
		v, multibyte, tail, err := strconv.UnquoteChar(lit[1:len(lit)-1], '\'')
		if err != nil {
			panic(err)
		}
		if tail != "" || multibyte {
			panic("invalid char: " + lit)
		}
		return byte(v)
	case "Text":
		return lit
	}
	panic("unsupported type: " + typ)
}

// -----------------------------------------------------------------------------

type Options struct {
	Scanner  tpl.Tokener
	ScanMode tpl.ScanMode
	GetValue func(typ string, lit string) interface{}
}

type Engine struct {
	*tpl.CompileRet
	doc Node
	stk *Node
}

var (
	optionsDef Options
)

var (
	InsertSemis = &Options{ScanMode: tpl.InsertSemis}
)

func New(grammar interface{}, options *Options) (p *Engine, err error) {

	if options == nil {
		options = &optionsDef
	}

	p = new(Engine)
	compiler := new(tpl.Compiler)

	initer := func() {
		p.doc.clear()
		p.stk = &p.doc
	}
	getValue := options.GetValue
	if getValue == nil {
		getValue = GetValue
	}
	marker := func(g tpl.Grammar, mark string) tpl.Grammar {
		markOld := mark
		isArray := false
		if strings.HasSuffix(mark, "s") {
			mark, isArray = mark[:len(mark)-1], true
		} else if strings.HasSuffix(mark, "Array") {
			mark, isArray = mark[:len(mark)-5], true
		}
		pos := strings.LastIndex(mark, "_")
		if pos < 0 {
			panic("invalid mark: " + markOld)
		}
		typ := mark[pos+1:]
		mark = mark[:pos]
		if typ != "Node" {
			return tpl.Action(g, func(tokens []tpl.Token, g tpl.Grammar) {
				p.stk.insertLeaf(mark, getValue(typ, tokens[0].Literal), isArray)
			})
		}
		return tpl.Transaction(g, func() interface{} {
			stk := p.stk
			trans := &transCtx{stk: stk, hd: stk.hd}
			p.stk = stk.insertNode(mark, isArray)
			return trans
		}, func(ctx interface{}, err error) {
			trans := ctx.(*transCtx)
			p.stk = trans.stk
			if err != nil {
				p.stk.hd = trans.hd
			}
		})
	}

	switch g := grammar.(type) {
	case string:
		compiler.Grammar = []byte(g)
	case []byte:
		compiler.Grammar = g
	default:
		panic("invalid grammar type")
	}
	compiler.Scanner = options.Scanner
	compiler.ScanMode = options.ScanMode
	compiler.Marker = marker
	compiler.Init = initer
	ret, err := compiler.Cl()
	if err != nil {
		return
	}

	p.CompileRet = &ret
	return
}

type transCtx struct {
	stk *Node
	hd  *Cons
}

// -----------------------------------------------------------------------------

func (p *Engine) Result() map[string]interface{} {

	return p.doc.toMap()
}

func (p *Engine) Parse(code []byte, fname string) (doc map[string]interface{}, err error) {

	err = p.MatchExactly(code, fname)
	if err != nil {
		return
	}
	return p.doc.toMap(), nil
}

func (p *Engine) SafeParse(code []byte, fname string) (doc map[string]interface{}, err error) {

	defer func() {
		if e := recover(); e != nil {
			switch v := e.(type) {
			case string:
				err = errors.New(v)
			case error:
				err = v
			default:
				panic(e)
			}
		}
	}()

	err = p.MatchExactly(code, fname)
	if err != nil {
		return
	}
	return p.doc.toMap(), nil
}

func (p *Engine) Extract(src string) (doc map[string]interface{}, err error) {

	return p.SafeParse([]byte(src), "")
}

func (p *Engine) SafeExtract(src string) (doc map[string]interface{}, err error) {

	return p.SafeParse([]byte(src), "")
}

// -----------------------------------------------------------------------------

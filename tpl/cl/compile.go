/*
 * Copyright (c) 2025 The XGo Authors (xgo.dev). All rights reserved.
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

package cl

import (
	"fmt"
	"os"
	"strconv"

	"github.com/goplus/xgo/tpl/ast"
	"github.com/goplus/xgo/tpl/matcher"
	"github.com/goplus/xgo/tpl/token"
	"github.com/qiniu/x/errors"
)

var (
	// ErrNoDocFound error
	ErrNoDocFound = errors.New("no document rule found")
)

// Result represents the result of compiling a set of rules.
type Result struct {
	Doc   *matcher.Var
	Rules map[string]*matcher.Var
}

type choice struct {
	m *matcher.Choices
	c *ast.Choice
}

type context struct {
	rules   map[string]*matcher.Var
	choices []choice
	errs    errors.List
	fset    *token.FileSet
}

func (p *context) newErrorf(pos token.Pos, format string, args ...any) error {
	return &matcher.Error{Fset: p.fset, Pos: pos, Msg: fmt.Sprintf(format, args...)}
}

func (p *context) addErrorf(pos token.Pos, format string, args ...any) {
	p.errs.Add(p.newErrorf(pos, format, args...))
}

func (p *context) addError(pos token.Pos, msg string) {
	p.errs.Add(&matcher.Error{Fset: p.fset, Pos: pos, Msg: msg})
}

// New compiles a set of rules from the given files.
func New(fset *token.FileSet, files ...*ast.File) (ret Result, err error) {
	return NewEx(nil, fset, files...)
}

// Config configures the behavior of the compiler.
type Config struct {
	RetProcs   map[string]any
	OnConflict func(fset *token.FileSet, c *ast.Choice, firsts [][]any, i, at int)
}

// NewEx compiles a set of rules from the given files.
func NewEx(conf *Config, fset *token.FileSet, files ...*ast.File) (ret Result, err error) {
	if conf == nil {
		conf = &Config{}
	}
	retProcs := conf.RetProcs
	rules := make(map[string]*matcher.Var)
	ctx := &context{rules: rules, fset: fset}
	for _, f := range files {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.Rule:
				ident := decl.Name
				name := ident.Name
				if old, ok := rules[name]; ok {
					oldPos := fset.Position(old.Pos)
					ctx.addErrorf(ident.Pos(),
						"duplicate rule `%s`, previous declaration at %v", name, oldPos)
					continue
				}
				v := matcher.NewVar(ident.Pos(), name)
				rules[name] = v
			default:
				ctx.addError(decl.Pos(), "unknown declaration")
			}
		}
	}
	var doc *matcher.Var
	for _, f := range files {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.Rule:
				ident := decl.Name
				name := ident.Name
				v := rules[name]
				if r, ok := compileExpr(decl.Expr, ctx); ok {
					v.RetProc = retProcs[name]
					if e := v.Assign(r); e != nil {
						ctx.addError(ident.Pos(), e.Error())
					}
					if doc == nil {
						doc = v
					}
				}
			}
		}
	}
	if doc == nil {
		err = ErrNoDocFound
		return
	}
	defer func() {
		if e := recover(); e != nil {
			switch e := e.(type) {
			case matcher.RecursiveError:
				ctx.addError(e.Pos, e.Error())
			default:
				panic(e)
			}
		}
		err = ctx.errs.ToError()
	}()
	onConflict := conf.OnConflict
	if onConflict == nil {
		onConflict = onConflictDefault
	}
	for _, item := range ctx.choices {
		item.m.CheckConflicts(func(firsts [][]any, i, at int) {
			onConflict(fset, item.c, firsts, i, at)
		})
	}
	ret = Result{doc, rules}
	return
}

func onConflictDefault(fset *token.FileSet, c *ast.Choice, firsts [][]any, i, at int) {
	pos := fset.Position(c.Options[i].Pos())
	LogConflict(pos, firsts, i, at)
}

// LogConflict logs a conflict between two choices.
func LogConflict(pos token.Position, firsts [][]any, i, at int) {
	fmt.Fprintf(os.Stderr, "%v: [WARN] conflict between %v and %v\n", pos, firsts[i], firsts[at])
}

var (
	idents = map[string]token.Token{
		"EOF":     token.EOF,
		"COMMENT": token.COMMENT,
		"IDENT":   token.IDENT,
		"INT":     token.INT,
		"FLOAT":   token.FLOAT,
		"IMAG":    token.IMAG,
		"CHAR":    token.CHAR,
		"STRING":  token.STRING,
		"RAT":     token.RAT,
		"UNIT":    token.UNIT,
		"LPAREN":  token.LPAREN,
		"RPAREN":  token.RPAREN,
		"LBRACK":  token.LBRACK,
		"RBRACK":  token.RBRACK,
		"LBRACE":  token.LBRACE,
		"RBRACE":  token.RBRACE,
	}
)

func compileExpr(expr ast.Expr, ctx *context) (matcher.Matcher, bool) {
	switch expr := expr.(type) {
	case *ast.Ident:
		name := expr.Name
		if v, ok := ctx.rules[name]; ok {
			return v, true
		} else if tok, ok := idents[name]; ok {
			return matcher.Token(tok), true
		}
		var quoteCh byte
		switch name {
		case "RAWSTRING":
			quoteCh = '`'
		case "QSTRING":
			quoteCh = '"'
		case "SPACE":
			return matcher.WhiteSpace(), true
		default:
			ctx.addErrorf(expr.Pos(), "`%s` is undefined", name)
		}
		return matcher.String(quoteCh), true
	case *ast.BasicLit:
		lit := expr.Value
		switch expr.Kind {
		case token.CHAR:
			v, multibyte, tail, e := strconv.UnquoteChar(lit[1:len(lit)-1], '\'')
			if e != nil {
				ctx.addErrorf(expr.Pos(), "invalid literal %s: %v", lit, e)
				break
			}
			if tail != "" || multibyte {
				ctx.addError(expr.Pos(), "invalid literal "+lit)
				break
			}
			return tokenExpr(token.Token(v), expr, ctx)
		case token.STRING:
			v, e := strconv.Unquote(lit)
			if e != nil {
				ctx.addError(expr.Pos(), "invalid literal "+lit)
				break
			}
			if v == "" {
				return matcher.True(), true
			}
			if c := v[0]; c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c == '_' {
				return matcher.Literal(token.IDENT, v), true
			}
			if t, ok := checkToken(v); ok {
				return tokenExpr(t, expr, ctx)
			}
			fallthrough
		default:
			ctx.addError(expr.Pos(), "invalid literal "+lit)
		}
	case *ast.Sequence:
		items := make([]matcher.Matcher, len(expr.Items))
		for i, item := range expr.Items {
			if r, ok := compileExpr(item, ctx); ok {
				items[i] = r
			} else {
				return nil, false
			}
		}
		return matcher.Sequence(items...), true
	case *ast.Choice:
		options := make([]matcher.Matcher, len(expr.Options))
		for i, option := range expr.Options {
			if r, ok := compileExpr(option, ctx); ok {
				options[i] = r
			} else {
				return nil, false
			}
		}
		ret := matcher.Choice(options...)
		ctx.choices = append(ctx.choices, choice{ret, expr})
		return ret, true
	case *ast.UnaryExpr:
		if x, ok := compileExpr(expr.X, ctx); ok {
			switch expr.Op {
			case token.QUESTION:
				return matcher.Repeat01(x), true
			case token.MUL:
				return matcher.Repeat0(x), true
			case token.ADD:
				return matcher.Repeat1(x), true
			default:
				ctx.addErrorf(expr.Pos(), "invalid token %v", expr.Op)
			}
		}
	case *ast.BinaryExpr:
		x, ok1 := compileExpr(expr.X, ctx)
		y, ok2 := compileExpr(expr.Y, ctx)
		if ok1 && ok2 {
			switch expr.Op {
			case token.REM: // %
				return matcher.List(x, y), true
			case token.INC: // ++
				return matcher.Adjoin(x, y), true
			default:
				ctx.addErrorf(expr.Pos(), "invalid token %v", expr.Op)
			}
		}
	default:
		ctx.addError(expr.Pos(), "unknown expression")
	}
	return nil, false
}

func tokenExpr(tok token.Token, expr *ast.BasicLit, ctx *context) (matcher.Matcher, bool) {
	if tok.Len() > 0 {
		return matcher.Token(tok), true
	}
	ctx.addErrorf(expr.Pos(), "invalid token: %s", expr.Value)
	return nil, false
}

func checkToken(v string) (ret token.Token, ok bool) {
	if len(v) == 1 {
		return token.Token(v[0]), true
	}
	token.ForEach(0, func(tok token.Token, lit string) int {
		if lit == v {
			ret, ok = tok, true
			return token.Break
		}
		return 0
	})
	return
}

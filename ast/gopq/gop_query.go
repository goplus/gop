/*
 Copyright 2020 The GoPlus Authors (goplus.org)

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package gopq

import (
	"errors"
	"os"
	"syscall"

	"github.com/goplus/gop/ast"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/token"
)

// -----------------------------------------------------------------------------

var (
	// ErrNotFound - not found
	ErrNotFound = syscall.ENOENT
	// ErrBreak - break
	ErrBreak = syscall.ELOOP
	// ErrTooManyNodes - too may nodes
	ErrTooManyNodes = errors.New("too many nodes")
)

// Node - node interface
type Node interface {
	ast.Node
	ForEach(filter func(node Node) error) error
	Obj() interface{}
}

// NodeEnum - node enumerator
type NodeEnum interface {
	ForEach(filter func(node Node) error) error
}

// NodeSet - node set
type NodeSet struct {
	Data NodeEnum
	Err  error
}

// NewSource calls ParseFile for all files with names ending in ".gop" in the
// directory specified by path and returns a map of package name -> package
// AST with all the packages found.
//
// If filter != nil, only the files with os.FileInfo entries passing through
// the filter (and ending in ".gop") are considered. The mode bits are passed
// to ParseFile unchanged. Position information is recorded in fset, which
// must not be nil.
//
// If the directory couldn't be read, a nil map and the respective error are
// returned. If a parse error occurred, a non-nil but incomplete map and the
// first error encountered are returned.
func NewSource(
	fset *token.FileSet, path string,
	filter func(os.FileInfo) bool, mode parser.Mode) (doc NodeSet, err error) {

	pkgs, err := parser.ParseDir(fset, path, filter, mode)
	if err != nil {
		return
	}
	return NodeSet{Data: &oneNode{astPackages(pkgs)}}, nil
}

// NewSourceFrom calls ParseFile for all files with names ending in ".gop" in the
// directory specified by path and returns a map of package name -> package
// AST with all the packages found.
//
// If filter != nil, only the files with os.FileInfo entries passing through
// the filter (and ending in ".gop") are considered. The mode bits are passed
// to ParseFile unchanged. Position information is recorded in fset, which
// must not be nil.
//
// If the directory couldn't be read, a nil map and the respective error are
// returned. If a parse error occurred, a non-nil but incomplete map and the
// first error encountered are returned.
func NewSourceFrom(
	fset *token.FileSet, fs parser.FileSystem, path string,
	filter func(os.FileInfo) bool, mode parser.Mode) (doc NodeSet, err error) {

	pkgs, err := parser.ParseFSDir(fset, fs, path, filter, mode)
	if err != nil {
		return
	}
	return NodeSet{Data: &oneNode{astPackages(pkgs)}}, nil
}

// Ok returns if node set is valid or not.
func (p NodeSet) Ok() bool {
	return p.Err == nil
}

// FuncDecl returns *ast.FuncDecl node set.
func (p NodeSet) FuncDecl() NodeSet {
	return p.Match(func(node Node) bool {
		_, ok := node.Obj().(*ast.FuncDecl)
		return ok
	})
}

// GenDecl returns *ast.GenDecl node set.
func (p NodeSet) GenDecl(tok token.Token) NodeSet {
	return p.Match(func(node Node) bool {
		decl, ok := node.Obj().(*ast.GenDecl)
		return ok && decl.Tok == tok
	})
}

// TypeSpec returns *ast.TypeSpec node set.
func (p NodeSet) TypeSpec() NodeSet {
	return p.GenDecl(token.TYPE).Child()
}

// VarSpec returns variables *ast.ValueSpec node set.
func (p NodeSet) VarSpec() NodeSet {
	return p.GenDecl(token.VAR).Child()
}

// ConstSpec returns constants *ast.ValueSpec node set.
func (p NodeSet) ConstSpec() NodeSet {
	return p.GenDecl(token.CONST).Child()
}

// ImportSpec returns *ast.ImportSpec node set.
func (p NodeSet) ImportSpec() NodeSet {
	return p.GenDecl(token.IMPORT).Child()
}

// -----------------------------------------------------------------------------

type oneNode struct {
	Node
}

func (p *oneNode) ForEach(filter func(node Node) error) error {
	return filter(p.Node)
}

func (p *oneNode) Cached() int {
	return 1
}

// One returns the first node as a node set.
func (p NodeSet) One() NodeSet {
	if _, ok := p.Data.(*oneNode); ok {
		return p
	}
	node, err := p.CollectOne()
	if err != nil {
		return NodeSet{Err: err}
	}
	return NodeSet{Data: &oneNode{node}}
}

// One creates a node set that only contains a signle node.
func One(node Node) NodeSet {
	return NodeSet{Data: &oneNode{node}}
}

// -----------------------------------------------------------------------------

type fixNodes struct {
	nodes []Node
}

func (p *fixNodes) ForEach(filter func(node Node) error) error {
	for _, node := range p.nodes {
		if filter(node) == ErrBreak {
			return ErrBreak
		}
	}
	return nil
}

func (p *fixNodes) Cached() int {
	return len(p.nodes)
}

// Nodes creates a fixed node set.
func Nodes(nodes ...Node) NodeSet {
	return NodeSet{Data: &fixNodes{nodes}}
}

// -----------------------------------------------------------------------------

type cached interface {
	Cached() int
}

// Cache caches node set.
func (p NodeSet) Cache() NodeSet {
	if _, ok := p.Data.(cached); ok {
		return p
	}
	nodes, err := p.Collect()
	if err != nil {
		return NodeSet{Err: err}
	}
	return NodeSet{Data: &fixNodes{nodes}}
}

// -----------------------------------------------------------------------------

type anyNodes struct {
	data NodeEnum
}

func (p *anyNodes) ForEach(filter func(node Node) error) error {
	return p.data.ForEach(func(node Node) error {
		return anyForEach(node, filter)
	})
}

func anyForEach(p Node, filter func(node Node) error) error {
	if err := filter(p); err == ErrBreak {
		return err
	}
	return p.ForEach(func(node Node) error {
		return anyForEach(node, filter)
	})
}

// Any returns deeply visiting node set.
func (p NodeSet) Any() (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &anyNodes{p.Data}}
}

// -----------------------------------------------------------------------------

type childNodes struct {
	data NodeEnum
}

func (p *childNodes) ForEach(filter func(node Node) error) error {
	return p.data.ForEach(func(node Node) error {
		return node.ForEach(filter)
	})
}

// Child returns child node set.
func (p NodeSet) Child() NodeSet {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &childNodes{p.Data}}
}

// -----------------------------------------------------------------------------

type matchedNodes struct {
	data  NodeEnum
	match func(node Node) bool
}

func (p *matchedNodes) ForEach(filter func(node Node) error) error {
	return p.data.ForEach(func(node Node) error {
		if p.match(node) {
			return filter(node)
		}
		return ErrNotFound
	})
}

// Match filters the node set.
func (p NodeSet) Match(match func(node Node) bool) (ret NodeSet) {
	if p.Err != nil {
		return p
	}
	return NodeSet{Data: &matchedNodes{p.Data, match}}
}

// -----------------------------------------------------------------------------

// Name returns names of the node set.
func (p NodeSet) Name() []string {
	return p.ToString(NameOf)
}

// ToString returns string values of the node set.
func (p NodeSet) ToString(str func(node Node) string) (items []string) {
	if p.Err != nil {
		return nil
	}
	p.Data.ForEach(func(node Node) error {
		items = append(items, str(node))
		return nil
	})
	return
}

// Collect collects all nodes of the node set.
func (p NodeSet) Collect() (items []Node, err error) {
	if p.Err != nil {
		return nil, p.Err
	}
	p.Data.ForEach(func(node Node) error {
		items = append(items, node)
		return nil
	})
	return
}

// CollectOne collects one node of a node set.
// If exactly is true, it returns ErrTooManyNodes when node set is more than one.
func (p NodeSet) CollectOne(exactly ...bool) (item Node, err error) {
	if p.Err != nil {
		return nil, p.Err
	}
	err = ErrNotFound
	if exactly != nil {
		if !exactly[0] {
			panic("please call `CollectOne()` instead of `CollectOne(false)`")
		}
		p.Data.ForEach(func(node Node) error {
			if err == ErrNotFound {
				item, err = node, nil
				return nil
			}
			err = ErrTooManyNodes
			return ErrBreak
		})
	} else {
		p.Data.ForEach(func(node Node) error {
			item, err = node, nil
			return ErrBreak
		})
	}
	return
}

// -----------------------------------------------------------------------------

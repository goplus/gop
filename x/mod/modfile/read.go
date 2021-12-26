/*
 * Copyright (c) 2021 The GoPlus Authors (goplus.org). All rights reserved.
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
package modfile

import (
	"golang.org/x/mod/modfile"
)

// A Position describes an arbitrary source position in a file, including the
// file, line, column, and byte offset.
type Position = modfile.Position

// An Expr represents an input element.
type Expr = modfile.Expr

// A Comment represents a single // comment.
type Comment = modfile.Comment

// Comments collects the comments associated with an expression.
type Comments = modfile.Comments

// A FileSyntax represents an entire gop.mod file.
type FileSyntax = modfile.FileSyntax

// A CommentBlock represents a top-level block of comments separate
// from any rule.
type CommentBlock = modfile.CommentBlock

// A Line is a single line of tokens.
type Line = modfile.Line

// A LineBlock is a factored block of lines, like
//
//	require (
//		"x"
//		"y"
//	)
//
type LineBlock = modfile.LineBlock

// An LParen represents the beginning of a parenthesized line block.
// It is a place to store suffix comments.
type LParen = modfile.LParen

// An RParen represents the end of a parenthesized line block.
// It is a place to store whole-line (before) comments.
type RParen = modfile.RParen

// ModulePath returns the module path from the gopmod file text.
// If it cannot find a module path, it returns an empty string.
// It is tolerant of unrelated problems in the gop.mod file.
func ModulePath(mod []byte) string {
	return modfile.ModulePath(mod)
}

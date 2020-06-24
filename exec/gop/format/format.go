/*
 Copyright 2020 damonchen (netubu@gmail.com)

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

package format

import (
	"github.com/qiniu/goplus/exec/gop/printer"
	"github.com/qiniu/goplus/parser"
	"github.com/qiniu/goplus/token"
)

var config = printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

// Source formats src in canonical gopfmt style and returns the result
// or an (I/O or syntax) error. src is expected to be a syntactically
// correct Gop source file, or a list of Gop declarations or statements.
//
// If src is a partial source file, the leading and trailing space of src
// is applied to the result (such that it has the same leading and trailing
// space as src), and the result is indented by the same amount as the first
// line of src containing code. Imports are not sorted for partial source files.
//
func Source(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return format(fset, file, src, config)
}

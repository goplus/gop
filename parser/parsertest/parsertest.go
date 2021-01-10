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

package parsertest

import (
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/goplus/gop/ast"
)

func sortedKeys(m interface{}) []string {
	iter := reflect.ValueOf(m).MapRange()
	keys := make([]string, 0, 8)
	for iter.Next() {
		key := iter.Key()
		keys = append(keys, key.String())
	}
	sort.Strings(keys)
	return keys
}

// PackageDebug prints debug information for a Go+ package.
func PackageDebug(w io.Writer, pkg *ast.Package) {
	fmt.Fprintf(w, "package %s\n", pkg.Name)
	paths := sortedKeys(pkg.Files)
	for _, path := range paths {
		fmt.Fprintf(w, "file %s\n", path)
		file := pkg.Files[path]
		_ = file
	}
}

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

package bytecode

import (
	"testing"

	"github.com/qiniu/x/ts"
)

// -----------------------------------------------------------------------------

func TestStack(test *testing.T) {
	t := ts.New(test)
	stk := NewStack()
	stk.Push(1)
	stk.Push(2)
	stk.Push(3)
	stk.Set(-1, 4)
	t.Call(stk.Get, -1).Equal(4)
	stk.SetLen(2)
	t.Call(stk.Get, -1).Equal(2)
}

// -----------------------------------------------------------------------------

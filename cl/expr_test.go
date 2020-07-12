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
package cl

import (
	"testing"
	"github.com/goplus/gop/ast/asttest"
)

var testDeleteClauses = map[string]testData{
	"delete_int_key": {`
					m:={1:1,2:2}
					delete(m,1)
					println(m)
					delete(m,3)
					println(m)
					delete(m,2)
					println(m)
					`, []string{"map[2:2]", "map[2:2]", "map[]"}},
	"delete_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					delete(m,"hi")
					println(m)
					delete(m,"Go+")
					println(m)
					`, []string{"map[Go+:2]", "map[Go+:2]", "map[]"}},
	"delete_var_string_key": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hello")
					println(m)
					a:="hi"
					delete(m,a)
					println(m)
					arr:=["Go+"]
					delete(m,arr[0])
					println(m)
					`, []string{"map[Go+:2]", "map[Go+:2]", "map[]"}},
	"delete_var_map_string_key": {`
					ma:=[{"hello":1,"Go+":2}]
					delete(ma[0],"hello")
					println(ma[0])
					a:="hi"
					delete(ma[0],a)
					println(ma[0])
					arr:=["Go+"]
					delete(ma[0],arr[0])
					println(ma[0])
					`, []string{"map[Go+:2]", "map[Go+:2]", "map[]"}},
	"delete_no_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m)
					`, []string{"_panic"}},
	"delete_multi_key_panic": {`
					m:={"hello":1,"Go+":2}
					delete(m,"hi","hi")
					`, []string{"_panic"}},
	"delete_not_map_panic": {`
					m:=[1,2,3]
					delete(m,1)
					`, []string{"_panic"}},
}

func TestDelete(t *testing.T) {
	for name, clause := range testDeleteClauses {
		testSingleStmt(name, t, asttest.NewSingleFileFS("/foo", "bar.gop", clause.clause), clause.wants)
	}
}

// -----------------------------------------------------------------------------

var testCopyClauses = map[string]testData{
	"copy_int": {`
					a:=[1,2,3]
					b:=[4,5,6]
					n:=copy(b,a)
					println(n)
					println(b)
					`, []string{"3", "[1 2 3]"}},
	"copy_string": {`
					a:=["hello"]
					b:=["hi"]
					n:=copy(b,a)
					println(n)
					println(b)
					`, []string{"1", "[hello]"}},
	"copy_byte_string": {`
					a:=[byte(65),byte(66),byte(67)]
					println(string(a))
					n:=copy(a,"abc")
					println(n)
					println(a)
					println(string(a))
					`, []string{"ABC", "3", "[97 98 99]", "abc"}},
	"copy_first_not_slice_panic": {`
					a:=1
					b:=[1,2,3]
					copy(a,b)
					println(a)
					`, []string{"_panic"}},
	"copy_second_not_slice_panic": {`
					a:=1
					b:=[1,2,3]
					copy(b,a)
					println(b)
					`, []string{"_panic"}},
	"copy_one_args_panic": {`
					a:=[1,2,3]
					copy(a)
					println(a)
					`, []string{"_panic"}},
	"copy_multi_args_panic": {`
					a:=[1,2,3]
					copy(a,a,a)
					println(a)
					`, []string{"_panic"}},
	"copy_string_panic": {`
					a:=[65,66,67]
					copy(a,"abc")
					println(a)
					`, []string{"_panic"}},
	"copy_different_type_panic": {`
					a:=[65,66,67]
					b:=[1.2,1.5,1.7]
					copy(b,a)
					copy(b,a)
					println(b)
					`, []string{"_panic"}},
	"copy_with_operation": {`
					a:=[65,66,67]
					b:=[1]
					println(copy(a,b)+copy(b,a)==2)
					`, []string{"true"}},
}

func TestCopy(t *testing.T) {
	for name, clause := range testCopyClauses {
		testSingleStmt(name, t, asttest.NewSingleFileFS("/foo", "bar.gop", clause.clause), clause.wants)
	}
}

// -----------------------------------------------------------------------------

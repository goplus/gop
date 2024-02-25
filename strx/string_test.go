/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
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

package strx

import (
	"testing"
)

func TestConcat(t *testing.T) {
	if ret := Concat("1", "23", "!"); ret != "123!" {
		t.Fatal("Concat:", ret)
	}
}

func TestBuild(t *testing.T) {
	if ret := NewBuilder(0).Build(); ret != "" {
		t.Fatal("NewBuilder(0):", ret)
	}
	if ret := NewBuilder(16).Add("1").AddByte('2', '3').AddByte('!').Build(); ret != "123!" {
		t.Fatal("TestBuild:", ret)
	}
}

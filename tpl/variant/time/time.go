/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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

package time

import (
	"time"

	"github.com/goplus/gop/tpl/variant"
)

// -----------------------------------------------------------------------------

func now(args ...any) any {
	if len(args) > 0 {
		panic("function call: arity mismatch")
	}
	return time.Now()
}

func init() {
	mod := variant.NewModule("time")
	mod.Insert("now", now)
}

// -----------------------------------------------------------------------------

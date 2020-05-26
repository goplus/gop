/*
 Copyright 2020 Qiniu Cloud (qiniu.com)

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

package exec

import "github.com/qiniu/x/log"

// -----------------------------------------------------------------------------

func checkPop(ctx *Context) interface{} {
	if ctx.Len() < 1 {
		log.Panicln("checkPop failed: no data.")
	}
	v := ctx.Get(-1)
	ctx.PopN(1)
	if ctx.Len() > 0 {
		log.Panicln("checkPop failed: too many data:", ctx.Len())
	}
	return v
}

// -----------------------------------------------------------------------------

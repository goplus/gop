/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package templ

import "testing"

const yapScriptIn = `
<html>
<body>
<pre>{{
$name := .Name
$base := .Base
range .Items
}}
<a href="{{.URL}}">{{.Name}}</a> (pv: {{.Meta.ViewCount}} nps: {{.NPS | printf "%.3f"}} by: {{.Meta.User}}) (<a href="/metadata/{{.ID}}">meta</a>){{
    if $base
        if .Ann}} (annotation: <a href="/annotation/{{.Ann}}?p=1">{{.Ann}}</a>){{
        end
    end
    with .Meta.Categories
        $n := len .
        if gt $n $base}} (under:{{
            range .
                if ne . $name}} <a href="/category/{{.}}?p=1">{{.}}</a>{{
                end
            end
            }}){{
        end
    end
end
}}
</pre>
</body>
</html>
`

const yapScriptOut = `
<html>
<body>
<pre>{{
$name := .Name}}{{
$base := .Base}}{{
range .Items
}}
<a href="{{.URL}}">{{.Name}}</a> (pv: {{.Meta.ViewCount}} nps: {{.NPS | printf "%.3f"}} by: {{.Meta.User}}) (<a href="/metadata/{{.ID}}">meta</a>){{
    if $base}}{{
        if .Ann}} (annotation: <a href="/annotation/{{.Ann}}?p=1">{{.Ann}}</a>){{
        end}}{{
    end}}{{
    with .Meta.Categories}}{{
        $n := len .}}{{
        if gt $n $base}} (under:{{
            range .}}{{
                if ne . $name}} <a href="/category/{{.}}?p=1">{{.}}</a>{{
                end}}{{
            end
            }}){{
        end}}{{
    end}}{{
end
}}
</pre>
</body>
</html>
`

func TestTranslate(t *testing.T) {
	if ret := Translate(yapScriptIn); ret != yapScriptOut {
		t.Fatal("TestTranslate:", len(ret), len(yapScriptOut), len(yapScriptIn)+11*4, ret)
	}
	noScript := "abc"
	if Translate(noScript) != noScript {
		t.Fatal("translate(noScript)")
	}
	noScriptEnd := `{{abc
efg
`
	noScriptEndOut := `{{abc}}{{
efg
`
	if ret := Translate(noScriptEnd); ret != noScriptEndOut {
		t.Fatal("translate(noScriptEnd):", ret)
	}
}

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
package env

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"

	"github.com/goplus/gop/cmd/internal/base"
	"github.com/goplus/gop/env"
)

// Cmd - gop env
var Cmd = &base.Command{
	UsageLine: "gop env [-json] [var ...]",
	Short:     "prints Go+ environment information.",
}

var (
	flag    = &Cmd.Flag
	envJson = flag.Bool("json", false, "prints Go environment information.")
)

func init() {
	Cmd.Run = runCmd
}

func runCmd(_ *base.Command, args []string) {
	err := flag.Parse(args)
	if err != nil {
		log.Fatalln("parse input arguments failed:", err)
	}

	var stdout bytes.Buffer

	cmd := exec.Command("go", "env", "-json")
	cmd.Env = os.Environ()
	cmd.Stdout = &stdout

	err = cmd.Run()
	if err != nil {
		log.Fatalln("run go env failed:", err)
	}

	var gopEnv map[string]interface{}
	if err := json.Unmarshal(stdout.Bytes(), &gopEnv); err != nil {
		log.Fatal("decode json of go env failed:", err)
	}

	gopEnv["BUILDDATE"] = env.BuildDate()
	gopEnv["BUILDREV"] = env.BuildRevision()
	gopEnv["GOPVERSION"] = env.Version()
	gopEnv["GOPROOT"] = env.GOPROOT()
	gopEnv["GOPATH"] = env.GOPATH()
	gopEnv["GOMODCACHE"] = env.GOMODCACHE()
	gopEnv["GOPMOD"], _ = env.GOPMOD("")
	gopEnv["HOME"] = env.HOME()

	onlyValues := false
	vars := flag.Args()
	if len(vars) > 0 {
		onlyValues = true
		newEnv := make(map[string]interface{})
		for _, v := range vars {
			if value, ok := gopEnv[v]; ok {
				newEnv[v] = value
			} else {
				newEnv[v] = ""
			}
		}
		gopEnv = newEnv
	}

	outputEnvVars(gopEnv, onlyValues, *envJson)
}

func outputEnvVars(gopEnv map[string]interface{}, onlyValues bool, outputJson bool) {
	if outputJson {
		b, err := json.Marshal(gopEnv)
		if err != nil {
			log.Fatal("encode json of go env failed:", err)
		}

		var out bytes.Buffer
		json.Indent(&out, b, "", "  ")
		fmt.Println(out.String())
	} else {
		keys := make([]string, 0, len(gopEnv))
		for k := range gopEnv {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := gopEnv[k]
			if onlyValues {
				fmt.Printf("%v\n", v)
			} else {
				fmt.Printf("%s=\"%v\"\n", k, v)
			}
		}
	}
}

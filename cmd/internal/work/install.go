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

package work

import (
	"fmt"
	"os"
	"os/exec"
)

func GoInstall(dir string, args ...string) error {
	gobin, err := exec.LookPath("go")
	if err != nil {
		return err
	}
	cmd := exec.Command(gobin, append([]string{gobin, "install"}, args...)...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("GOBIN=%v", GOPBIN()))
	return cmd.Run()
}

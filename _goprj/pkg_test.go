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

package goprj_test

import (
	"testing"

	"github.com/qiniu/goplus/goprj"
	"github.com/qiniu/x/log"
)

func init() {
	log.SetFlags(log.Llevel)
	log.SetOutputLevel(log.Ldebug)
}

func Test(t *testing.T) {
	prj := goprj.NewProject()
	pkg, err := prj.OpenPackage(".")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.Source().Name != "goprj" {
		t.Fatal("please run test in this package directory")
	}
	if pkg.ThisModule().PkgPath() != "github.com/qiniu/goplus/goprj" {
		t.Fatal("PkgPath:", pkg.ThisModule().PkgPath())
	}
	return
	log.Debug("------------------------------------------------------")
	pkg2, err := pkg.LoadPackage("github.com/qiniu/goplus/modutil")
	if err != nil {
		t.Fatal(err)
	}
	if pkg2.Source().Name != "modutil" {
		t.Fatal("please run test in this package directory")
	}
	log.Debug("------------------------------------------------------")
	prjPath := "github.com/visualfc/fastmod"
	pkg3, err := pkg2.LoadPackage(prjPath)
	if err != nil {
		t.Fatal(err)
	}
	if pkg3.ThisModule().VersionPkgPath() != "github.com/visualfc/fastmod@v1.3.3" {
		t.Fatal("PkgPath:", pkg3.ThisModule().VersionPkgPath())
	}
}

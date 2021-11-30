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

package cl_test

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/gox"
)

func newTwoFileFS(dir string, fname, data string, fname2 string, data2 string) *parsertest.MemFS {
	return parsertest.NewMemFS(map[string][]string{
		dir: {fname, fname2},
	}, map[string]string{
		path.Join(dir, fname):  data,
		path.Join(dir, fname2): data2,
	})
}

func init() {
	cl.RegisterClassFileType(".tgmx", ".tspx", "github.com/goplus/gop/cl/internal/spx", "math")
}

func gopSpxTest(t *testing.T, gmx, spxcode, expected string) {
	gopSpxTestEx(t, gmx, spxcode, expected, "index.tgmx", "bar.tspx")
}

func gopSpxTestEx(t *testing.T, gmx, spxcode, expected, gmxfile, spxfile string) {
	cl.SetDisableRecover(true)
	defer cl.SetDisableRecover(false)

	fs := newTwoFileFS("/foo", spxfile, spxcode, gmxfile, gmx)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	conf := *baseConf.Ensure()
	bar := pkgs["main"]
	pkg, err := cl.NewPackage("", bar, &conf)
	if err != nil {
		t.Fatal("NewPackage:", err)
	}
	var b bytes.Buffer
	err = gox.WriteTo(&b, pkg, false)
	if err != nil {
		t.Fatal("gox.WriteTo failed:", err)
	}
	result := b.String()
	if result != expected {
		t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
	}
}

func gopSpxErrorTestEx(t *testing.T, msg, gmx, spxcode, gmxfile, spxfile string) {
	fs := newTwoFileFS("/foo", spxfile, spxcode, gmxfile, gmx)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", nil, 0)
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	conf := *baseConf.Ensure()
	conf.NoFileLine = false
	conf.WorkingDir = "/foo"
	conf.TargetDir = "/foo"
	bar := pkgs["main"]
	_, err = cl.NewPackage("", bar, &conf)
	if err == nil {
		t.Fatal("no error?")
	}
	if ret := err.Error(); ret != msg {
		t.Fatalf("\nError: \"%s\"\nExpected: \"%s\"\n", ret, msg)
	}
}

func TestSpxError(t *testing.T) {
	gopSpxErrorTestEx(t, `./Game.tgmx:4:2: cannot assign value to field in class file`, `
var (
	Kai Kai
	userScore int = 100
)
`, `
println "hi"
`, "Game.tgmx", "Kai.tspx")

	gopSpxErrorTestEx(t, `./Kai.tspx:3:2: missing field type in class file`, `
var (
	Kai Kai
	userScore int
)
`, `
var (
	id = 100
)
println "hi"
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxBasic(t *testing.T) {
	gopSpxTestEx(t, `
func onInit() {
	for {
	}
}
`, `
func onMsg(msg string) {
	for {
		say "Hi"
	}
}
`, `package main

import spx "github.com/goplus/gop/cl/internal/spx"

type Game struct {
	*spx.MyGame
}

func (this *Game) onInit() {
	for {
		spx.SchedNow()
	}
}

type Kai struct {
	spx.Sprite
	*Game
}

func (this *Kai) onMsg(msg string) {
	for {
		spx.Sched()
		this.Say("Hi")
	}
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxBasic2(t *testing.T) {
	gopSpxTest(t, `
import (
	"fmt"
)

const (
	Foo = 1
)

func bar() {
}

func onInit() {
	Foo
	bar
	fmt.Println("Hi")
}
`, ``, `package main

import (
	fmt "fmt"
	spx "github.com/goplus/gop/cl/internal/spx"
)

const Foo = 1

type index struct {
	*spx.MyGame
}

func (this *index) bar() {
}
func (this *index) onInit() {
	Foo
	this.bar()
	fmt.Println("Hi")
}
`)
}

func TestSpxMethod(t *testing.T) {
	gopSpxTestEx(t, `
func onInit() {
	sched
	broadcast "msg1"
	testIntValue = 1
	x := round(1.2)
}
`, `
func onInit() {
	setCostume "kai-a"
	play "recordingWhere"
	say "Where do you come from?", 2
	broadcast "msg2"
}
`, `package main

import (
	spx "github.com/goplus/gop/cl/internal/spx"
	math "math"
)

type Game struct {
	*spx.MyGame
}

func (this *Game) onInit() {
	spx.Sched()
	this.Broadcast__0("msg1")
	spx.TestIntValue = 1
	x := math.Round(1.2)
}

type bar struct {
	spx.Sprite
	*Game
}

func (this *bar) onInit() {
	this.SetCostume("kai-a")
	this.Play("recordingWhere")
	this.Say("Where do you come from?", 2)
	this.Broadcast__0("msg2")
}
`, "Game.tgmx", "bar.tspx")
}

func TestSpxVar(t *testing.T) {
	gopSpxTestEx(t, `
var (
	Kai Kai
)

func onInit() {
	Kai.clone()
	broadcast("msg1")
}
`, `
var (
	a int
)

func onInit() {
	a = 1
}

func onCloned() {
	say("Hi")
}
`, `package main

import spx "github.com/goplus/gop/cl/internal/spx"

type Game struct {
	*spx.MyGame
	Kai Kai
}
type Kai struct {
	spx.Sprite
	*Game
	a int
}

func (this *Game) onInit() {
	this.Kai.Clone()
	this.Broadcast__0("msg1")
}
func (this *Kai) onInit() {
	this.a = 1
}
func (this *Kai) onCloned() {
	this.Say("Hi")
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxRun(t *testing.T) {
	gopSpxTestEx(t, `
var (
	Kai Kai
	t   Sound
)

run "hzip://open.qiniu.us/weather/res.zip"
`, `
println "Hi"
`, `package main

import (
	fmt "fmt"
	spx "github.com/goplus/gop/cl/internal/spx"
)

type index struct {
	*spx.MyGame
	Kai Kai
	t   spx.Sound
}
type Kai struct {
	spx.Sprite
	*index
}

func (this *index) MainEntry() {
	spx.Gopt_MyGame_Run(this, "hzip://open.qiniu.us/weather/res.zip")
}
func main() {
	spx.Gopt_MyGame_Main(new(index))
}
func (this *Kai) Main() {
	fmt.Println("Hi")
}
`, "index.tgmx", "Kai.tspx")
}

func TestSpx2(t *testing.T) {
	gopSpxTestEx(t, `
println("Hi")
`, `
func onMsg(msg string) {
}
`, `package main

import (
	fmt "fmt"
	spx2 "github.com/goplus/gop/cl/internal/spx2"
)

type Game struct {
	spx2.Game
}

func (this *Game) MainEntry() {
	fmt.Println("Hi")
}
func main() {
	new(Game).Main()
}

type Kai struct {
	spx2.Sprite
	*Game
}

func (this *Kai) onMsg(msg string) {
}
`, "Game.t2gmx", "Kai.t2spx")
}

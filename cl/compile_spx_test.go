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
	"testing"

	"github.com/goplus/gop/cl"
	"github.com/goplus/gop/parser"
	"github.com/goplus/gop/parser/parsertest"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/mod/gopmod"
)

func lookupClass(ext string) (c *gopmod.Project, ok bool) {
	switch ext {
	case ".tgmx", ".tspx":
		return &gopmod.Project{
			Ext: ".tgmx", Class: "*MyGame",
			Works:    []*gopmod.Class{{Ext: ".tspx", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx", "math"}}, true
	case ".t2gmx", ".t2spx", ".t2spx2":
		return &gopmod.Project{
			Ext: ".t2gmx", Class: "Game",
			Works: []*gopmod.Class{{Ext: ".t2spx", Class: "Sprite"},
				{Ext: ".t2spx2", Class: "Sprite2"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	case ".t3spx", ".t3spx2":
		return &gopmod.Project{
			Works: []*gopmod.Class{{Ext: ".t3spx", Class: "Sprite"},
				{Ext: ".t3spx2", Class: "Sprite2"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	}
	return
}

func spxParserConf() parser.Config {
	return parser.Config{
		IsClass: func(ext string) (isProj bool, ok bool) {
			c, ok := lookupClass(ext)
			if ok {
				isProj = (c.Ext == ext)
			}
			return
		},
	}
}

func gopSpxTest(t *testing.T, gmx, spxcode, expected string) {
	gopSpxTestEx(t, gmx, spxcode, expected, "index.tgmx", "bar.tspx")
}

func gopSpxTestEx(t *testing.T, gmx, spxcode, expected, gmxfile, spxfile string) {
	gopSpxTestExConf(t, "gopSpxTest", gblConf, gmx, spxcode, expected, gmxfile, spxfile)
}

func gopSpxTestExConf(t *testing.T, name string, conf *cl.Config, gmx, spxcode, expected, gmxfile, spxfile string) {
	t.Run(name, func(t *testing.T) {
		cl.SetDisableRecover(true)
		defer cl.SetDisableRecover(false)

		fs := parsertest.NewTwoFilesFS("/foo", spxfile, spxcode, gmxfile, gmx)
		pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", spxParserConf())
		if err != nil {
			scanner.PrintError(os.Stderr, err)
			t.Fatal("ParseFSDir:", err)
		}
		bar := pkgs["main"]
		pkg, err := cl.NewPackage("", bar, conf)
		if err != nil {
			t.Fatal("NewPackage:", err)
		}
		var b bytes.Buffer
		err = pkg.WriteTo(&b)
		if err != nil {
			t.Fatal("gox.WriteTo failed:", err)
		}
		result := b.String()
		if result != expected {
			t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
		}
	})
}

func gopSpxErrorTestEx(t *testing.T, msg, gmx, spxcode, gmxfile, spxfile string) {
	fs := parsertest.NewTwoFilesFS("/foo", spxfile, spxcode, gmxfile, gmx)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", spxParserConf())
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	conf := *gblConf
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
	gopSpxErrorTestEx(t, `./Game.tgmx:6:2: userScore redeclared
	./Game.tgmx:5:2 other declaration of userScore`, `
import "bytes"
var (
	Kai Kai
	userScore int
	userScore string
)
`, `
println "hi"
`, "Game.tgmx", "Kai.tspx")

	gopSpxErrorTestEx(t, `./Kai.tspx:4:2: id redeclared
	./Kai.tspx:3:2 other declaration of id`, `
var (
	Kai Kai
	userScore int
)
`, `
var (
	id int
	id string
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
	TestIntValue = 1
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

func (this *Game) onInit() {
	this.Kai.Clone()
	this.Broadcast__0("msg1")
}

type Kai struct {
	spx.Sprite
	*Game
	a int
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

var x float64 = rand(1.2)

run "hzip://open.qiniu.us/weather/res.zip"
`, `
println "Hi"
`, `package main

import (
	fmt "fmt"
	spx "github.com/goplus/gop/cl/internal/spx"
)

var x float64 = spx.Rand__1(1.2)

type index struct {
	*spx.MyGame
	Kai Kai
	t   spx.Sound
}

func (this *index) MainEntry() {
	spx.Gopt_MyGame_Run(this, "hzip://open.qiniu.us/weather/res.zip")
}
func main() {
	spx.Gopt_MyGame_Main(new(index))
}

type Kai struct {
	spx.Sprite
	*index
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

	gopSpxTestEx(t, `
println("Hi, Sprite2")
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
	fmt.Println("Hi, Sprite2")
}
func main() {
	new(Game).Main()
}

type Kai struct {
	spx2.Sprite2
	*Game
}

func (this *Kai) onMsg(msg string) {
}
`, "Game.t2gmx", "Kai.t2spx2")

	gopSpxTestEx(t, `
println("Hi, Sprite")
`, `
func onMsg(msg string) {
}
`, `package main

import (
	fmt "fmt"
	spx2 "github.com/goplus/gop/cl/internal/spx2"
)

type Dog struct {
	spx2.Sprite
}

func (this *Dog) Main() {
	fmt.Println("Hi, Sprite")
}

type Kai struct {
	spx2.Sprite2
}

func (this *Kai) onMsg(msg string) {
}
`, "Dog.t3spx", "Kai.t3spx2")
}

func TestSpxMainEntry(t *testing.T) {
	conf := *gblConf
	conf.Importer = nil
	conf.NoAutoGenMain = false

	gopSpxTestExConf(t, "Nocode", &conf, `
`, `
`, `package main

import spx2 "github.com/goplus/gop/cl/internal/spx2"

type Game struct {
	spx2.Game
}

func (this *Game) MainEntry() {
}
func main() {
	new(Game).Main()
}
`, "Game.t2gmx", "Kai.t2spx")
	gopSpxTestExConf(t, "OnlyGmx", &conf, `
var (
	Kai Kai
)
`, `
`, `package main

import spx2 "github.com/goplus/gop/cl/internal/spx2"

type Game struct {
	spx2.Game
	Kai Kai
}

func (this *Game) MainEntry() {
}
func main() {
	new(Game).Main()
}

type Kai struct {
	spx2.Sprite
	*Game
}
`, "Game.t2gmx", "Kai.t2spx")

	gopSpxTestExConf(t, "KaiAndGmx", &conf, `
var (
	Kai Kai
)
func MainEntry() {
	println "Hi"
}
`, `
func Main() {
	println "Hello"
}
func onMsg(msg string) {
}
`, `package main

import (
	fmt "fmt"
	spx2 "github.com/goplus/gop/cl/internal/spx2"
)

type Game struct {
	spx2.Game
	Kai Kai
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

func (this *Kai) Main() {
	fmt.Println("Hello")
}
func (this *Kai) onMsg(msg string) {
}
`, "Game.t2gmx", "Kai.t2spx")
}

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
	"github.com/goplus/gop/parser/fsx/memfs"
	"github.com/goplus/gop/scanner"
	"github.com/goplus/mod/modfile"
)

func lookupClass(ext string) (c *modfile.Project, ok bool) {
	switch ext {
	case ".tgmx", ".tspx":
		return &modfile.Project{
			Ext: ".tgmx", Class: "*MyGame",
			Works:    []*modfile.Class{{Ext: ".tspx", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx", "math"}}, true
	case ".t2gmx", ".t2spx", ".t2spx2":
		return &modfile.Project{
			Ext: ".t2gmx", Class: "Game",
			Works: []*modfile.Class{{Ext: ".t2spx", Class: "Sprite"},
				{Ext: ".t2spx2", Class: "Sprite2"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	case "_t3spx.gox", ".t3spx2":
		return &modfile.Project{
			Works: []*modfile.Class{{Ext: "_t3spx.gox", Class: "Sprite"},
				{Ext: ".t3spx2", Class: "Sprite2"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx2"}}, true
	case "_spx.gox":
		return &modfile.Project{
			Ext: "_spx.gox", Class: "Game",
			Works:    []*modfile.Class{{Ext: "_spx.gox", Class: "Sprite"}},
			PkgPaths: []string{"github.com/goplus/gop/cl/internal/spx3", "math"},
			Import:   []*modfile.Import{{Path: "github.com/goplus/gop/cl/internal/spx3/jwt"}}}, true
	case "_xtest.gox":
		return &modfile.Project{
			Ext: "_xtest.gox", Class: "App",
			Works:    []*modfile.Class{{Ext: "_xtest.gox", Class: "Case"}},
			PkgPaths: []string{"github.com/goplus/gop/test", "testing"}}, true
	}
	return
}

func spxParserConf() parser.Config {
	return parser.Config{
		ClassKind: func(fname string) (isProj bool, ok bool) {
			ext := modfile.ClassExt(fname)
			c, ok := lookupClass(ext)
			if ok {
				isProj = c.IsProj(ext, fname)
			}
			return
		},
	}
}

func gopSpxTest(t *testing.T, gmx, spxcode, expected string) {
	gopSpxTestEx(t, gmx, spxcode, expected, "index.tgmx", "bar.tspx")
}

func gopSpxTestEx(t *testing.T, gmx, spxcode, expected, gmxfile, spxfile string) {
	gopSpxTestExConf(t, "gopSpxTest", gblConf, gmx, spxcode, expected, gmxfile, spxfile, "")
}

func gopSpxTestEx2(t *testing.T, gmx, spxcode, expected, gmxfile, spxfile, resultFile string) {
	gopSpxTestExConf(t, "gopSpxTest", gblConf, gmx, spxcode, expected, gmxfile, spxfile, resultFile)
}

func gopSpxTestExConf(t *testing.T, name string, conf *cl.Config, gmx, spxcode, expected, gmxfile, spxfile, resultFile string) {
	t.Run(name, func(t *testing.T) {
		cl.SetDisableRecover(true)
		defer cl.SetDisableRecover(false)

		fs := memfs.TwoFiles("/foo", spxfile, spxcode, gmxfile, gmx)
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
		err = pkg.WriteTo(&b, resultFile)
		if err != nil {
			t.Fatal("gogen.WriteTo failed:", err)
		}
		result := b.String()
		if result != expected {
			t.Fatalf("\nResult:\n%s\nExpected:\n%s\n", result, expected)
		}
	})
}

func gopSpxErrorTestEx(t *testing.T, msg, gmx, spxcode, gmxfile, spxfile string) {
	fs := memfs.TwoFiles("/foo", spxfile, spxcode, gmxfile, gmx)
	pkgs, err := parser.ParseFSDir(gblFset, fs, "/foo", spxParserConf())
	if err != nil {
		scanner.PrintError(os.Stderr, err)
		t.Fatal("ParseFSDir:", err)
	}
	conf := *gblConf
	conf.RelativeBase = "/foo"
	conf.Recorder = nil
	conf.NoFileLine = false
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
	gopSpxErrorTestEx(t, `Game.tgmx:6:2: userScore redeclared
	Game.tgmx:5:2 other declaration of userScore`, `
import "bytes"
var (
	Kai Kai
	userScore int
	userScore string
)
`, `
println "hi"
`, "Game.tgmx", "Kai.tspx")

	gopSpxErrorTestEx(t, `Kai.tspx:4:2: id redeclared
	Kai.tspx:3:2 other declaration of id`, `
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

initGameApp
`, `
func onMsg(msg string) {
	for {
		say "Hi"
	}
}
`, `package main

import "github.com/goplus/gop/cl/internal/spx"

type Game struct {
	*spx.MyGame
}
type Kai struct {
	spx.Sprite
	*Game
}

func (this *Game) onInit() {
	for {
		spx.SchedNow()
	}
}
func (this *Game) MainEntry() {
	this.InitGameApp()
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onMsg(msg string) {
	for {
		spx.Sched()
		this.Say("Hi")
	}
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
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
	"fmt"
	"github.com/goplus/gop/cl/internal/spx"
)

const Foo = 1

type bar struct {
	spx.Sprite
	*index
}
type index struct {
	*spx.MyGame
}

func (this *index) bar() {
}
func (this *index) onInit() {
	this.bar()
	fmt.Println("Hi")
}
func (this *index) MainEntry() {
}
func (this *index) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *bar) Classfname() string {
	return "bar"
}
func (this *bar) Main() {
}
func main() {
	new(index).Main()
}
`)
}

func TestEnvOp(t *testing.T) {
	gopSpxTest(t, `
echo ${PATH}, $id
`, ``, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx"
)

type bar struct {
	spx.Sprite
	*index
}
type index struct {
	*spx.MyGame
}

func (this *index) MainEntry() {
	fmt.Println(this.Gop_Env("PATH"), this.Gop_Env("id"))
}
func (this *index) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *bar) Classfname() string {
	return "bar"
}
func (this *bar) Main() {
}
func main() {
	new(index).Main()
}
`)
}

func TestSpxGopEnv(t *testing.T) {
	gopSpxTest(t, `
echo "${PATH}"
`, ``, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx"
	"strconv"
)

type bar struct {
	spx.Sprite
	*index
}
type index struct {
	*spx.MyGame
}

func (this *index) MainEntry() {
	fmt.Println(strconv.Itoa(this.Gop_Env("PATH")))
}
func (this *index) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *bar) Classfname() string {
	return "bar"
}
func (this *bar) Main() {
}
func main() {
	new(index).Main()
}
`)
}

func TestSpxGopExec(t *testing.T) {
	gopSpxTest(t, `
vim "a.txt"
vim
ls 10
capout => { ls }
capout => { ls "-l" }
`, ``, `package main

import "github.com/goplus/gop/cl/internal/spx"

type bar struct {
	spx.Sprite
	*index
}
type index struct {
	*spx.MyGame
}

func (this *index) MainEntry() {
	this.Gop_Exec("vim", "a.txt")
	this.Gop_Exec("vim")
	this.Ls(10)
	this.Capout(func() {
		this.Gop_Exec("ls")
	})
	this.Capout(func() {
		this.Gop_Exec("ls", "-l")
	})
}
func (this *index) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *bar) Classfname() string {
	return "bar"
}
func (this *bar) Main() {
}
func main() {
	new(index).Main()
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
	"github.com/goplus/gop/cl/internal/spx"
	"math"
)

type Game struct {
	*spx.MyGame
}
type bar struct {
	spx.Sprite
	*Game
}

func (this *Game) onInit() {
	spx.Sched()
	this.Broadcast__0("msg1")
	spx.TestIntValue = 1
	x := math.Round(1.2)
}
func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *bar) onInit() {
	this.SetCostume("kai-a")
	this.Play("recordingWhere")
	this.Say("Where do you come from?", 2)
	this.Broadcast__0("msg2")
}
func (this *bar) Classfname() string {
	return "bar"
}
func (this *bar) Main() {
}
func main() {
	new(Game).Main()
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

import "github.com/goplus/gop/cl/internal/spx"

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
	spx.Gopt_Sprite_Clone__0(this.Kai)
	this.Broadcast__0("msg1")
}
func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onInit() {
	this.a = 1
}
func (this *Kai) onCloned() {
	this.Say("Hi")
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
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
	"fmt"
	"github.com/goplus/gop/cl/internal/spx"
)

type Kai struct {
	spx.Sprite
	*index
}
type index struct {
	*spx.MyGame
	Kai Kai
	t   spx.Sound
}

var x float64 = spx.Rand__1(1.2)

func (this *index) MainEntry() {
	spx.Gopt_MyGame_Run(this, "hzip://open.qiniu.us/weather/res.zip")
}
func (this *index) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) Main() {
	fmt.Println("Hi")
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func main() {
	new(index).Main()
}
`, "index.tgmx", "Kai.tspx")
}

func TestSpxRunWithWorkers(t *testing.T) {
	gopSpxTestEx(t, `
var (
	Kai Kai
)

run
`, `
echo jwt.token("Hi")
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx3"
	"github.com/goplus/gop/cl/internal/spx3/jwt"
)

type Kai struct {
	spx3.Sprite
	*Game
}
type Game struct {
	spx3.Game
	Kai Kai
}

func (this *Game) MainEntry() {
	this.Run()
}
func (this *Game) Main() {
	spx3.Gopt_Game_Main(this, new(Kai))
}
func (this *Kai) Main(_gop_arg0 string) {
	this.Sprite.Main(_gop_arg0)
	fmt.Println(jwt.Token("Hi"))
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func main() {
	new(Game).Main()
}
`, "main_spx.gox", "Kai_spx.gox")
}

func TestSpxNewObj(t *testing.T) {
	gopSpxTestEx(t, `
a := new
a.run
b := new(Sprite)
echo b.name
`, ``, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx3"
)

type Kai struct {
	spx3.Sprite
	*Game
}
type Game struct {
	spx3.Game
}

func (this *Game) MainEntry() {
	a := spx3.New()
	a.Run()
	b := new(spx3.Sprite)
	fmt.Println(b.Name())
}
func (this *Game) Main() {
	spx3.Gopt_Game_Main(this, new(Kai))
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main(_gop_arg0 string) {
	this.Sprite.Main(_gop_arg0)
}
func main() {
	new(Game).Main()
}
`, "main_spx.gox", "Kai_spx.gox")
}

func TestSpx2(t *testing.T) {
	gopSpxTestEx(t, `
println("Hi")
`, `
func onMsg(msg string) {
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx2"
)

type Game struct {
	spx2.Game
}
type Kai struct {
	spx2.Sprite
	*Game
}

func (this *Game) MainEntry() {
	fmt.Println("Hi")
}
func (this *Game) Main() {
	(*spx2.Game).Main(&this.Game)
}
func (this *Kai) onMsg(msg string) {
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.t2gmx", "Kai.t2spx")
}

func TestSpx3(t *testing.T) {
	gopSpxTestEx(t, `
println("Hi, Sprite2")
`, `
func onMsg(msg string) {
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx2"
)

type Game struct {
	spx2.Game
}
type Kai struct {
	spx2.Sprite2
	*Game
}

func (this *Game) MainEntry() {
	fmt.Println("Hi, Sprite2")
}
func (this *Game) Main() {
	(*spx2.Game).Main(&this.Game)
}
func (this *Kai) onMsg(msg string) {
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.t2gmx", "Kai.t2spx2")
}

func TestSpx4(t *testing.T) {
	gopSpxTestEx(t, `
println("Hi, Sprite")
`, `
func onMsg(msg string) {
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx2"
)

type Dog struct {
	spx2.Sprite
}
type Kai struct {
	spx2.Sprite2
}

func (this *Dog) Main() {
	fmt.Println("Hi, Sprite")
}
func (this *Dog) Classfname() string {
	return "Dog"
}
func (this *Kai) onMsg(msg string) {
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
`, "Dog_t3spx.gox", "Kai.t3spx2")
}

func TestSpxMainEntry(t *testing.T) {
	conf := *gblConf
	conf.Importer = nil
	conf.NoAutoGenMain = false

	gopSpxTestExConf(t, "Nocode", &conf, `
`, `
`, `package main

import "github.com/goplus/gop/cl/internal/spx2"

type Game struct {
	spx2.Game
}
type Kai struct {
	spx2.Sprite
	*Game
}

func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	(*spx2.Game).Main(&this.Game)
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.t2gmx", "Kai.t2spx", "")
	gopSpxTestExConf(t, "OnlyGmx", &conf, `
var (
	Kai Kai
)
`, `
`, `package main

import "github.com/goplus/gop/cl/internal/spx2"

type Game struct {
	spx2.Game
	Kai Kai
}
type Kai struct {
	spx2.Sprite
	*Game
}

func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	(*spx2.Game).Main(&this.Game)
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.t2gmx", "Kai.t2spx", "")

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
	"fmt"
	"github.com/goplus/gop/cl/internal/spx2"
)

type Game struct {
	spx2.Game
	Kai Kai
}
type Kai struct {
	spx2.Sprite
	*Game
}

func (this *Game) MainEntry() {
	fmt.Println("Hi")
}
func (this *Game) Main() {
	(*spx2.Game).Main(&this.Game)
}
func (this *Kai) Main() {
	fmt.Println("Hello")
}
func (this *Kai) onMsg(msg string) {
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func main() {
	new(Game).Main()
}
`, "Game.t2gmx", "Kai.t2spx", "")
}

func TestSpxGoxBasic(t *testing.T) {
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

import "github.com/goplus/gop/cl/internal/spx"

type Game struct {
	*spx.MyGame
}
type Kai struct {
	spx.Sprite
	*Game
}

func (this *Game) onInit() {
	for {
		spx.SchedNow()
	}
}
func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onMsg(msg string) {
	for {
		spx.Sched()
		this.Say("Hi")
	}
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxClone(t *testing.T) {
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

type info struct {
	x int
	y int
}

func onInit() {
	a = 1
	clone
	clone info{1,2}
	clone &info{1,2}
}

func onCloned() {
	say("Hi")
}
`, `package main

import "github.com/goplus/gop/cl/internal/spx"

type info struct {
	x int
	y int
}
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
	spx.Gopt_Sprite_Clone__0(this.Kai)
	this.Broadcast__0("msg1")
}
func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onInit() {
	this.a = 1
	spx.Gopt_Sprite_Clone__0(this)
	spx.Gopt_Sprite_Clone__1(this, info{1, 2})
	spx.Gopt_Sprite_Clone__1(this, &info{1, 2})
}
func (this *Kai) onCloned() {
	this.Say("Hi")
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxErrorSel(t *testing.T) {
	gopSpxErrorTestEx(t, `Kai.tspx:2:9: this.pos undefined (type *Kai has no field or method pos)`, `
println "hi"
`, `
println this.pos
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxMethodSel(t *testing.T) {
	gopSpxTestEx(t, `
sendMessage "Hi"
`, `
func onMsg(msg string) {
}
`, `package main

import "github.com/goplus/gop/cl/internal/spx"

type Game struct {
	*spx.MyGame
}
type Kai struct {
	spx.Sprite
	*Game
}

func (this *Game) MainEntry() {
	this.SendMessage("Hi")
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onMsg(msg string) {
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxPkgOverload(t *testing.T) {
	gopSpxTestEx(t, `
println "Hi"
`, `
func onMsg(msg string) {
	this.position.add 100,200
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx"
)

type Game struct {
	*spx.MyGame
}
type Kai struct {
	spx.Sprite
	*Game
}

func (this *Game) MainEntry() {
	fmt.Println("Hi")
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onMsg(msg string) {
	this.Position().Add__0(100, 200)
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxSelection(t *testing.T) {
	gopSpxTestEx(t, `
println "hi"
`, `
import "fmt"
func onMsg(msg string) {
	fmt.println msg
	this.position.add 100,200
	position.add 100,200
	position.X += 100
	println position.X
	this.vector.add 100,200
	vector.add 100,200
	vector.X += 100
	vector.self.X += 100
	vector.self.Y += 200
	vector.self.add position.X,position.Y
	println vector.X
	println vector.self.self
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/cl/internal/spx"
)

type Game struct {
	*spx.MyGame
}
type Kai struct {
	spx.Sprite
	*Game
}

func (this *Game) MainEntry() {
	fmt.Println("hi")
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (this *Kai) onMsg(msg string) {
	fmt.Println(msg)
	this.Position().Add__0(100, 200)
	this.Position().Add__0(100, 200)
	this.Position().X += 100
	fmt.Println(this.Position().X)
	this.Vector().Add__0(100, 200)
	this.Vector().Add__0(100, 200)
	this.Vector().X += 100
	this.Vector().Self().X += 100
	this.Vector().Self().Y += 200
	this.Vector().Self().Add__0(this.Position().X, this.Position().Y)
	fmt.Println(this.Vector().X)
	fmt.Println(this.Vector().Self().Self())
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Main() {
}
func main() {
	new(Game).Main()
}
`, "Game.tgmx", "Kai.tspx")
}

func TestSpxOverload(t *testing.T) {
	gopSpxTestEx(t, `
var (
	Kai Kai
)

func onInit() {
	Kai.onKey "hello", key => {
	}
}
`, `
var (
	a int
)

type Mesh struct {
}

func (p *Mesh) Name() string {
	return "hello"
}

var (
	m1 = &Mesh{}
	m2 = &Mesh{}
)

onKey "hello", => {
}
onKey "hello", key => {
}
onKey ["1"], => {
}
onKey ["2"], key => {
}
onKey [m1, m2], => {
}
onKey [m1, m2], key => {
}
onKey ["a"], ["b"], key => {
}
onKey ["a"], [m1, m2], key => {
}
onKey ["a"], nil, key => {
}
onKey 100, 200
onKey2 "hello", key => {
}
`, `package main

import "github.com/goplus/gop/cl/internal/spx"

type Mesh struct {
}
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
	spx.Gopt_Sprite_OnKey__1(this.Kai, "hello", func(key string) {
	})
}
func (this *Game) MainEntry() {
}
func (this *Game) Main() {
	spx.Gopt_MyGame_Main(this)
}
func (p *Mesh) Name() string {
	return "hello"
}

var m1 = &Mesh{}
var m2 = &Mesh{}

func (this *Kai) Main() {
	spx.Gopt_Sprite_OnKey__0(this, "hello", func() {
	})
	spx.Gopt_Sprite_OnKey__1(this, "hello", func(key string) {
	})
	spx.Gopt_Sprite_OnKey__2(this, []string{"1"}, func() {
	})
	spx.Gopt_Sprite_OnKey__3(this, []string{"2"}, func(key string) {
	})
	spx.Gopt_Sprite_OnKey__4(this, []spx.Mesher{m1, m2}, func() {
	})
	spx.Gopt_Sprite_OnKey__5(this, []spx.Mesher{m1, m2}, func(key spx.Mesher) {
	})
	spx.Gopt_Sprite_OnKey__6(this, []string{"a"}, []string{"b"}, func(key string) {
	})
	spx.Gopt_Sprite_OnKey__7(this, []string{"a"}, []spx.Mesher{m1, m2}, func(key string) {
	})
	spx.Gopt_Sprite_OnKey__6(this, []string{"a"}, nil, func(key string) {
	})
	spx.Gopt_Sprite_OnKey__8(this, 100, 200)
	spx.Gopt_Sprite_OnKey2(this, "hello", func(key string) {
	})
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func main() {
	new(Game).Main()
}
`, "Game.tgmx", "Kai.tspx")
}

func TestTestClassFile(t *testing.T) {
	gopSpxTestEx2(t, `
println "Hi"
`, `
t.log "Hi"
t.run "a test", t => {
	t.fatal "failed"
}
`, `package main

import (
	"fmt"
	"github.com/goplus/gop/test"
	"testing"
)

type caseFoo struct {
	test.Case
	*App
}
type App struct {
	test.App
}

func (this *App) MainEntry() {
	fmt.Println("Hi")
}
func (this *caseFoo) Main() {
	this.T().Log("Hi")
	this.T().Run("a test", func(t *testing.T) {
		t.Fatal("failed")
	})
}
func (this *caseFoo) Classfname() string {
	return "Foo"
}
func TestFoo(t *testing.T) {
	test.Gopt_Case_TestMain(new(caseFoo), t)
}
func TestMain(m *testing.M) {
	test.Gopt_App_TestMain(new(App), m)
}
`, "main_xtest.gox", "Foo_xtest.gox", "_test")
}

func TestTestClassFile2(t *testing.T) {
	gopSpxTestEx2(t, `
println "Hi"
`, `
t.log "Hi"
`, `package main

import (
	"github.com/goplus/gop/test"
	"testing"
)

type case_foo struct {
	test.Case
}

func (this *case_foo) Main() {
	this.T().Log("Hi")
}
func (this *case_foo) Classfname() string {
	return "foo"
}
func Test_foo(t *testing.T) {
	test.Gopt_Case_TestMain(new(case_foo), t)
}
`, "main.gox", "foo_xtest.gox", "_test")
}

func TestGopxNoFunc(t *testing.T) {
	gopClTestFile(t, `
var (
	a int
)
`, `package main

type foo struct {
	a int
}
`, "foo.gox")
}

func TestClassFileGopx(t *testing.T) {
	gopClTestFile(t, `
var (
	BaseClass
	Width, Height float64
	*AggClass
)

type BaseClass struct{
	x int
	y int
}
type AggClass struct{}

func Area() float64 {
	return Width * Height
}
`, `package main

type BaseClass struct {
	x int
	y int
}
type AggClass struct {
}
type Rect struct {
	BaseClass
	Width  float64
	Height float64
	*AggClass
}

func (this *Rect) Area() float64 {
	return this.Width * this.Height
}
`, "Rect.gox")
	gopClTestFile(t, `
import "bytes"
var (
	bytes.Buffer
)
func test(){}
`, `package main

import "bytes"

type Rect struct {
	bytes.Buffer
}

func (this *Rect) test() {
}
`, "Rect.gox")
	gopClTestFile(t, `
import "bytes"
var (
	*bytes.Buffer
)
func test(){}
`, `package main

import "bytes"

type Rect struct {
	*bytes.Buffer
}

func (this *Rect) test() {
}
`, "Rect.gox")
	gopClTestFile(t, `
import "bytes"
var (
	*bytes.Buffer "spec:\"buffer\""
	a int "json:\"a\""
	b int
)
func test(){}
`, `package main

import "bytes"

type Rect struct {
	*bytes.Buffer `+"`spec:\"buffer\"`"+`
	a             int `+"`json:\"a\"`"+`
	b             int
}

func (this *Rect) test() {
}
`, "Rect.gox")
}

func TestClassFileMember(t *testing.T) {
	gopClTestFile(t, `type Engine struct {
}

func (e *Engine) EnterPointerLock() {
}

func (e *Engine) SetEnable(b bool) {
}

func Engine() *Engine {
	return &Engine{}
}

func Test() {
	engine.setEnable true
	engine.enterPointerLock
}
`, `package main

type Engine struct {
}
type Rect struct {
}

func (e *Engine) EnterPointerLock() {
}
func (e *Engine) SetEnable(b bool) {
}
func (this *Rect) Engine() *Engine {
	return &Engine{}
}
func (this *Rect) Test() {
	this.Engine().SetEnable(true)
	this.Engine().EnterPointerLock()
}
`, "Rect.gox")
}

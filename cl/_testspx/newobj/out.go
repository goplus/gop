package main

import (
	"fmt"

	"github.com/goplus/xgo/cl/internal/spx3"
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
	_gop_obj0 := &Kai{Game: this}
	spx3.Gopt_Game_Main(this, _gop_obj0)
}
func (this *Kai) Main(_gop_arg0 string) {
	this.Sprite.Main(_gop_arg0)
}
func (this *Kai) Classfname() string {
	return "Kai"
}
func (this *Kai) Classclone() spx3.Handler {
	_gop_ret := *this
	return &_gop_ret
}
func main() {
	new(Game).Main()
}

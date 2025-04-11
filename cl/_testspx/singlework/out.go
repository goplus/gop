package main

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
func (this *Kai) Classclone() spx3.Handler {
	_gop_ret := *this
	return &_gop_ret
}
func main() {
	new(Game).Main()
}

package main

import "github.com/goplus/xgo/cl/internal/mcp"

type foo struct {
	mcp.Prompt
	*Game
}
type Tool_hello struct {
	mcp.Tool
	*Game
}
type Game struct {
	mcp.Game
	foo *foo
}

func (this *Game) MainEntry() {
	this.Server("protos")
}
func (this *Game) Main() {
	_gop_obj0 := &Tool_hello{Game: this}
	_gop_lst1 := []mcp.ToolProto{_gop_obj0}
	_gop_obj1 := &foo{Game: this}
	this.foo = _gop_obj1
	_gop_lst2 := []mcp.PromptProto{_gop_obj1}
	mcp.Gopt_Game_Main(this, nil, _gop_lst1, _gop_lst2)
}
func (this *foo) Main(_gop_arg0 *mcp.Tool) string {
	this.Prompt.Main(_gop_arg0)
	return "Hi"
}
func (this *Tool_hello) Main(_gop_arg0 string) int {
	this.Tool.Main(_gop_arg0)
	return -1
}
func main() {
	new(Game).Main()
}

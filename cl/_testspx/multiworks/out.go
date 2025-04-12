package main

import "github.com/goplus/gop/cl/internal/mcp"

type foo struct {
	mcp.Prompt
	*Game
}
type hello struct {
	mcp.Tool
	*Game
}
type Game struct {
	mcp.Game
}

func (this *Game) MainEntry() {
	this.Server("protos")
}
func (this *Game) Main() {
	mcp.Gopt_Game_Main(this, nil, []mcp.ToolProto{new(hello)}, []mcp.PromptProto{new(foo)})
}
func (this *foo) Main(_gop_arg0 *mcp.Tool) string {
	this.Prompt.Main(_gop_arg0)
	return "Hi"
}
func (this *hello) Main(_gop_arg0 string) int {
	this.Tool.Main(_gop_arg0)
	return -1
}
func main() {
	new(Game).Main()
}

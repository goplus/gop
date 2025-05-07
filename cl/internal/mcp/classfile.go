package mcp

const (
	GopPackage = true
)

type Game struct {
}

func New() *Game {
	return nil
}

func (p *Game) initGame() {}

func (p *Game) Server(name string) {}

type Tool struct {
}

func (p *Tool) Main(name string) int {
	return 0
}

type Prompt struct {
}

func (p *Prompt) Main(*Tool) string {
	return ""
}

type Resource struct {
}

func (p *Resource) Main() {
}

type ToolProto interface {
	Main(name string) int
}

type PromptProto interface {
	Main(*Tool) string
}

type ResourceProto interface {
	Main()
}

func Gopt_Game_Main(game interface{ initGame() }, resources []ResourceProto, tools []ToolProto, prompts []PromptProto) {
}

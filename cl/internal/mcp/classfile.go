/*
 * Copyright (c) 2025 The GoPlus Authors (goplus.org). All rights reserved.
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

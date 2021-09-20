/*
 Copyright 2021 The GoPlus Authors (goplus.org)

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

package spx

const (
	GopPackage = true

	Gop_game   = "Game"
	Gop_title  = "on"
	Gop_params = "onMsg(, _gop_data interface{}); onCloned(_gop_data interface{})"
)

type Sound string

type Game struct {
}

func (p *Game) Broadcast__0(msg string) {
}

func (p *Game) Broadcast__1(msg string, wait bool) {
}

func (p *Game) Broadcast__2(msg string, data interface{}, wait bool) {
}

func (p *Game) Play(media string, wait ...bool) {
}

func Run(game interface{}, resource string) error {
	return nil
}

func Sched() {
}

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
	Gop_game   = "MyGame"
	Gop_sched  = "Sched,SchedNow"
)

type Sound string

type MyGame struct {
}

func (p *MyGame) Broadcast__0(msg string) {
}

func (p *MyGame) Broadcast__1(msg string, wait bool) {
}

func (p *MyGame) Broadcast__2(msg string, data interface{}, wait bool) {
}

func (p *MyGame) Play(media string, wait ...bool) {
}

func Gopt_MyGame_Run(game interface{}, resource string) error {
	return nil
}

func Sched() {
}

func SchedNow() {
}

var (
	TestIntValue int
)

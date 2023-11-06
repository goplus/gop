/*
 * Copyright (c) 2023 The GoPlus Authors (goplus.org). All rights reserved.
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

package jsonrpc2test

import (
	"context"
	"io"
	"net"
	"testing"
)

func TestNetPipeDone(t *testing.T) {
	np := &netPiper{
		done:   make(chan struct{}, 1),
		dialed: make(chan io.ReadWriteCloser),
	}
	np.done <- struct{}{}
	if f, err := np.Accept(context.Background()); f != nil || err != net.ErrClosed {
		t.Fatal("np.Accept:", f, err)
	}
	np.done <- struct{}{}
	if f, err := np.Dial(context.Background()); f != nil || err != net.ErrClosed {
		t.Fatal("np.Dial:", f, err)
	}
}

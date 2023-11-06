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

package stdio

import (
	"context"
	"net"
	"testing"

	"github.com/goplus/gop/x/jsonrpc2"
)

func TestConn(t *testing.T) {
	jsonrpc2.SetDebug(jsonrpc2.DbgFlagAll)
	c := new(conn)
	c.Dial(context.Background())
	if err := c.Close(); err != nil {
		t.Fatal("Close failed:", err)
	}
	if err := c.Close(); err != net.ErrClosed {
		t.Fatal("Close:", err)
	}
	if _, err := c.Read(nil); err != net.ErrClosed {
		t.Fatal("Read:", err)
	}
	if _, err := c.Write(nil); err != net.ErrClosed {
		t.Fatal("Write:", err)
	}
}

func TestListener(t *testing.T) {
	l := new(listener)
	if l.Dialer() != nil {
		t.Fatal("Dialer: not nil?")
	}
	c, err := l.Dial(context.Background())
	if err != nil {
		t.Fatal("Dial failed:", err)
	}
	if _, err := l.Dial(context.Background()); err != ErrTooManyConnections {
		t.Fatal("Dial:", err)
	}
	if err := c.Close(); err != nil {
		t.Fatal("Close failed:", err)
	}
	if err := c.Close(); err != net.ErrClosed {
		t.Fatal("Close:", err)
	}
	if _, err := c.Read(nil); err != net.ErrClosed {
		t.Fatal("Read:", err)
	}
	if _, err := c.Write(nil); err != net.ErrClosed {
		t.Fatal("Write:", err)
	}
	l.Close()
	if err := l.Close(); err != net.ErrClosed {
		t.Fatal("l.Close:", err)
	}
}

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
	"io"
	"os"
	"sync/atomic"

	"github.com/goplus/gop/x/fakenet"
	"github.com/goplus/gop/x/jsonrpc2"
	"github.com/goplus/gop/x/jsonrpc2/jsonrpc2test"
)

// -----------------------------------------------------------------------------

type listener struct {
}

func (p *listener) Close() error {
	return nil
}

// Accept blocks waiting for an incoming connection to the listener.
func (p *listener) Accept(context.Context) (io.ReadWriteCloser, error) {
	connCnt := &connCnt[server]
	if atomic.AddInt32(connCnt, 1) != 1 {
		atomic.AddInt32(connCnt, -1)
		return nil, ErrTooManyConnections
	}
	return fakenet.NewConn("stdio", os.Stdin, os.Stdout), nil
}

func (p *listener) Dialer() jsonrpc2.Dialer {
	return nil
}

// Listener returns a jsonrpc2.Listener based on stdin and stdout.
func Listener(allowLocalDail bool) jsonrpc2.Listener {
	if allowLocalDail {
		return jsonrpc2test.NetPipeListener()
	}
	return &listener{}
}

// -----------------------------------------------------------------------------

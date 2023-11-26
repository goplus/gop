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
	"errors"
	"io"
	"sync/atomic"

	"github.com/goplus/gop/x/fakenet"
	"github.com/goplus/gop/x/jsonrpc2"
)

var (
	ErrTooManyConnections = errors.New("too many connections")
)

const (
	client = iota
	server
)

var (
	connCnt [2]int32
)

// -----------------------------------------------------------------------------

type dialer struct {
	in  io.ReadCloser
	out io.WriteCloser
}

// Dial returns a new communication byte stream to a listening server.
func (p *dialer) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	dailCnt := &connCnt[client]
	if atomic.AddInt32(dailCnt, 1) != 1 {
		atomic.AddInt32(dailCnt, -1)
		return nil, ErrTooManyConnections
	}
	return fakenet.NewConn("stdio.dialer", p.in, p.out), nil
}

// Dialer returns a jsonrpc2.Dialer based on in and out.
func Dialer(in io.ReadCloser, out io.WriteCloser) jsonrpc2.Dialer {
	return &dialer{in: in, out: out}
}

// Dial makes a new connection based on in and out, wraps the returned
// reader and writer using the framer to make a stream, and then builds a
// connection on top of that stream using the binder.
//
// The returned Connection will operate independently using the Preempter and/or
// Handler provided by the Binder, and will release its own resources when the
// connection is broken, but the caller may Close it earlier to stop accepting
// (or sending) new requests.
func Dial(in io.ReadCloser, out io.WriteCloser, binder jsonrpc2.Binder, onDone func()) (*jsonrpc2.Connection, error) {
	return jsonrpc2.Dial(context.Background(), Dialer(in, out), binder, onDone)
}

// -----------------------------------------------------------------------------

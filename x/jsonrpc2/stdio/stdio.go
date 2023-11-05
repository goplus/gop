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
	"net"
	"os"
	"sync/atomic"

	"github.com/goplus/gop/x/jsonrpc2"
)

// -----------------------------------------------------------------------------

var (
	ErrTooManyConnections = errors.New("too many connections")
)

var (
	connCnt int32
)

type conn struct {
	closed bool
}

func (p *conn) Read(b []byte) (n int, err error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	return os.Stdin.Read(b)
}

func (p *conn) Write(b []byte) (n int, err error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	return os.Stdout.Write(b)
}

func (p *conn) Close() error {
	if p.closed {
		return net.ErrClosed
	}
	p.closed = true
	atomic.AddInt32(&connCnt, -1)
	return nil
}

// Dial returns a new communication byte stream to a listening server.
func (p *conn) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	if p.closed {
		return nil, net.ErrClosed
	}
	if atomic.AddInt32(&connCnt, 1) != 1 {
		atomic.AddInt32(&connCnt, -1)
		return nil, ErrTooManyConnections
	}
	return p, nil
}

// Accept blocks waiting for an incoming connection to the listener.
func (p *conn) Accept(context.Context) (io.ReadWriteCloser, error) {
	if p.closed {
		return nil, net.ErrClosed
	}
	if atomic.AddInt32(&connCnt, 1) != 1 {
		atomic.AddInt32(&connCnt, -1)
		return nil, ErrTooManyConnections
	}
	return p, nil
}

// -----------------------------------------------------------------------------

func (p *conn) Dialer() jsonrpc2.Dialer {
	return nil
}

// -----------------------------------------------------------------------------

// Dialer returns a jsonrpc2.Dialer based on stdin and stdout.
func Dialer() jsonrpc2.Dialer {
	return new(conn)
}

// Dial makes a new connection based on stdin and stdout, wraps the returned
// reader and writer using the framer to make a stream, and then builds a
// connection on top of that stream using the binder.
//
// The returned Connection will operate independently using the Preempter and/or
// Handler provided by the Binder, and will release its own resources when the
// connection is broken, but the caller may Close it earlier to stop accepting
// (or sending) new requests.
func Dial(binder jsonrpc2.Binder) (*jsonrpc2.Connection, error) {
	return jsonrpc2.Dial(context.Background(), Dialer(), binder)
}

// -----------------------------------------------------------------------------

// Listener returns a jsonrpc2.Listener based on stdin and stdout.
func Listener() jsonrpc2.Listener {
	return new(conn)
}

// NewServer starts a new server listening for incoming connections and returns
// it.
// This returns a fully running and connected server, it does not block on
// the listener.
// You can call Wait to block on the server, or Shutdown to get the sever to
// terminate gracefully.
// To notice incoming connections, use an intercepting Binder.
func NewServer(binder jsonrpc2.Binder) *jsonrpc2.Server {
	return jsonrpc2.NewServer(context.Background(), Listener(), binder)
}

// -----------------------------------------------------------------------------

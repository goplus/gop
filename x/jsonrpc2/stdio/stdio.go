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
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"

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
	n, err = os.Stdin.Read(b)
	for n == 0 { // retry (to support localDailer)
		time.Sleep(time.Millisecond)
		n, err = os.Stdin.Read(b)
	}
	return
}

func (p *conn) Write(b []byte) (n int, err error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	return os.Stdout.Write(b)
}

func (p *conn) Close() error {
	if jsonrpc2.Verbose {
		log.Println("==> stdio.conn.Close")
	}
	if p.closed {
		return net.ErrClosed
	}
	p.closed = true
	atomic.AddInt32(&connCnt, -1)
	return nil
}

// Dial returns a new communication byte stream to a listening server.
func (p *conn) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	if atomic.AddInt32(&connCnt, 1) != 1 {
		atomic.AddInt32(&connCnt, -1)
		return nil, ErrTooManyConnections
	}
	return p, nil
}

// -----------------------------------------------------------------------------

type newConn struct {
	in, out *os.File
	r, w    *os.File
	oStdin  *os.File
	oStdout *os.File
	closed  bool
}

func (p *newConn) Read(b []byte) (n int, err error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	return p.in.Read(b)
}

func (p *newConn) Write(b []byte) (n int, err error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	return p.out.Write(b)
}

func (p *newConn) Close() error {
	if jsonrpc2.Verbose {
		log.Println("==> stdio.newConn.Close")
	}
	if p.closed {
		return net.ErrClosed
	}
	p.closed = true
	atomic.AddInt32(&dailCnt, -1)
	os.Stdin, os.Stdout = p.oStdin, p.oStdout
	p.in.Close()
	p.w.Close()
	p.r.Close()
	p.out.Close()
	return nil
}

// -----------------------------------------------------------------------------

type listener struct {
	closed    bool
	localDail bool
}

func (p *listener) Close() error {
	if jsonrpc2.Verbose {
		log.Println("==> stdio.listener.Close")
	}
	if p.closed {
		return net.ErrClosed
	}
	p.closed = true
	return nil
}

// Accept blocks waiting for an incoming connection to the listener.
func (p *listener) Accept(context.Context) (io.ReadWriteCloser, error) {
	if atomic.AddInt32(&connCnt, 1) != 1 {
		atomic.AddInt32(&connCnt, -1)
		return nil, ErrTooManyConnections
	}
	return new(conn), nil
}

var (
	dailCnt int32
)

// Dial returns a new communication byte stream to a listening server.
func (p *listener) Dial(ctx context.Context) (ret io.ReadWriteCloser, err error) {
	if atomic.AddInt32(&dailCnt, 1) != 1 {
		atomic.AddInt32(&dailCnt, -1)
		return nil, ErrTooManyConnections
	}
	in, w, err := os.Pipe()
	if err != nil {
		return
	}
	r, out, err := os.Pipe()
	if err != nil {
		in.Close()
		w.Close()
		return
	}
	oStdin, oStdout := os.Stdin, os.Stdout
	os.Stdin, os.Stdout = r, w
	return &newConn{in: in, out: out, r: r, w: w, oStdin: oStdin, oStdout: oStdout}, nil
}

func (p *listener) Dialer() jsonrpc2.Dialer {
	if p.localDail {
		return p
	}
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
func Listener(allowLocalDail bool) jsonrpc2.Listener {
	l := &listener{localDail: allowLocalDail}
	return l
}

// NewServer starts a new server listening for incoming connections and returns
// it.
// This returns a fully running and connected server, it does not block on
// the listener.
// You can call Wait to block on the server, or Shutdown to get the sever to
// terminate gracefully.
// To notice incoming connections, use an intercepting Binder.
func NewServer(binder jsonrpc2.Binder, allowLocalDail bool) *jsonrpc2.Server {
	return jsonrpc2.NewServer(context.Background(), Listener(allowLocalDail), binder)
}

// -----------------------------------------------------------------------------

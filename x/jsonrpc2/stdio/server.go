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

type serverConn struct {
	closed bool
}

func (p *serverConn) Read(b []byte) (n int, err error) {
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

func (p *serverConn) Write(b []byte) (n int, err error) {
	if p.closed {
		return 0, net.ErrClosed
	}
	return os.Stdout.Write(b)
}

func (p *serverConn) Close() error {
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
	return new(serverConn), nil
}

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

// Listener returns a jsonrpc2.Listener based on stdin and stdout.
func Listener(allowLocalDail bool) jsonrpc2.Listener {
	l := &listener{localDail: allowLocalDail}
	return l
}

// -----------------------------------------------------------------------------

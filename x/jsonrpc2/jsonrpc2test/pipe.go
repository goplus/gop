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

// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonrpc2test

import (
	"context"
	"io"
	"net"

	"github.com/goplus/gop/x/jsonrpc2"
)

// NetPipeListener returns a new Listener that listens using net.Pipe.
// It is only possibly to connect to it using the Dialer returned by the
// Dialer method, each call to that method will generate a new pipe the other
// side of which will be returned from the Accept call.
func NetPipeListener() jsonrpc2.Listener {
	return &netPiper{
		done:   make(chan struct{}),
		dialed: make(chan io.ReadWriteCloser),
	}
}

// netPiper is the implementation of Listener build on top of net.Pipes.
type netPiper struct {
	done   chan struct{}
	dialed chan io.ReadWriteCloser
}

// Accept blocks waiting for an incoming connection to the listener.
func (l *netPiper) Accept(context.Context) (io.ReadWriteCloser, error) {
	// Block until the pipe is dialed or the listener is closed,
	// preferring the latter if already closed at the start of Accept.
	select {
	case <-l.done:
		return nil, net.ErrClosed
	default:
	}
	select {
	case rwc := <-l.dialed:
		return rwc, nil
	case <-l.done:
		return nil, net.ErrClosed
	}
}

// Close will cause the listener to stop listening. It will not close any connections that have
// already been accepted.
func (l *netPiper) Close() error {
	// unblock any accept calls that are pending
	close(l.done)
	return nil
}

func (l *netPiper) Dialer() jsonrpc2.Dialer {
	return l
}

func (l *netPiper) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	client, server := net.Pipe()

	select {
	case l.dialed <- server:
		return client, nil

	case <-l.done:
		client.Close()
		server.Close()
		return nil, net.ErrClosed
	}
}

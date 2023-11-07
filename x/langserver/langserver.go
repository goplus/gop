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

package langserver

import (
	"context"
	"encoding/json"
	"log"

	"github.com/goplus/gop"
	"github.com/goplus/gop/x/gopprojs"
	"github.com/goplus/gop/x/jsonrpc2"
)

// -----------------------------------------------------------------------------

// Listener is implemented by protocols to accept new inbound connections.
type Listener = jsonrpc2.Listener

// Server is a running server that is accepting incoming connections.
type Server = jsonrpc2.Server

// Config holds the options for new connections.
type Config struct {
	// Framer allows control over the message framing and encoding.
	// If nil, HeaderFramer will be used.
	Framer jsonrpc2.Framer
}

func NewServer(ctx context.Context, listener Listener, conf *Config) (ret *Server) {
	h := new(handler)
	ret = jsonrpc2.NewServer(ctx, listener, jsonrpc2.BinderFunc(
		func(ctx context.Context, c *jsonrpc2.Connection) (ret jsonrpc2.ConnectionOptions) {
			if conf != nil {
				ret.Framer = conf.Framer
			}
			ret.Handler = h
			ret.OnInternalError = h.OnInternalError
			return
		}))
	h.server = ret
	return
}

// -----------------------------------------------------------------------------

type handler struct {
	server *Server
}

func (p *handler) OnInternalError(err error) {
	panic("jsonrpc2: " + err.Error())
}

func (p *handler) Changed(files []string) {
	log.Println("Changed:", files)
}

func (p *handler) Handle(ctx context.Context, req *jsonrpc2.Request) (result interface{}, err error) {
	switch req.Method {
	case methodChanged:
		var files []string
		err = json.Unmarshal(req.Params, &files)
		if err != nil {
			return
		}
		p.Changed(files)
	case methodGenGo:
		var pattern []string
		err = json.Unmarshal(req.Params, &pattern)
		if err != nil {
			return
		}
		err = GenGo(pattern...)
	case methodShutdown:
		p.server.Shutdown()
	}
	return
}

func GenGo(pattern ...string) (err error) {
	projs, err := gopprojs.ParseAll(pattern...)
	if err != nil {
		return
	}
	for _, proj := range projs {
		switch v := proj.(type) {
		case *gopprojs.DirProj:
			gop.GenGoEx(v.Dir, nil, true, 0)
		case *gopprojs.PkgPathProj:
			if v.Path == "builtin" {
				continue
			}
			gop.GenGoPkgPathEx("", v.Path, nil, true, 0)
		}
	}
	return
}

// -----------------------------------------------------------------------------

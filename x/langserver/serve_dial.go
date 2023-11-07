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
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/goplus/gop/x/jsonrpc2/stdio"
)

func fatal(err error) {
	log.Fatalln(err)
}

// -----------------------------------------------------------------------------

// ServeAndDialConfig represents the configuration of ServeAndDial.
type ServeAndDialConfig struct {
	// LogFile is where the LangServer application log saves to (optional).
	// Default is ~/.gop/serve.log
	LogFile string

	// OnError is to customize how to process errors (optional).
	// It should panic in any case.
	OnError func(err error)
}

// ServeAndDial executes a command as a LangServer, makes a new connection to it
// and returns a client of the LangServer based on the connection.
func ServeAndDial(conf *ServeAndDialConfig, cmd string, args ...string) Client {
	if conf == nil {
		conf = new(ServeAndDialConfig)
	}
	onErr := conf.OnError
	if onErr == nil {
		onErr = fatal
	}
	logFile := conf.LogFile
	if logFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			onErr(err)
		}
		gopDir := home + "/.gop"
		err = os.MkdirAll(gopDir, 0755)
		if err != nil {
			onErr(err)
		}
		logFile = gopDir + "/serve.log"
	}

	in, w := io.Pipe()
	r, out := io.Pipe()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer r.Close()
		defer w.Close()

		f, err := os.Create(logFile)
		if err != nil {
			onErr(err)
		}
		defer f.Close()

		cmd := exec.CommandContext(ctx, cmd, args...)
		cmd.Stdin = r
		cmd.Stdout = w
		cmd.Stderr = f
		cmd.Run()
	}()

	c, err := Open(ctx, stdio.Dialer(in, out), cancel)
	if err != nil {
		onErr(err)
	}
	return c
}

// -----------------------------------------------------------------------------

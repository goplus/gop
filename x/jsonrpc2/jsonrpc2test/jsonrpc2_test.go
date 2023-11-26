// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonrpc2test_test

import (
	"context"
	"testing"

	"github.com/goplus/gop/x/jsonrpc2"
	"github.com/goplus/gop/x/jsonrpc2/jsonrpc2test"
	"github.com/goplus/gop/x/jsonrpc2/jsonrpc2test/cases"
)

func TestNetPipe(t *testing.T) {
	jsonrpc2.SetDebug(jsonrpc2.DbgFlagCall)
	ctx := context.Background()
	listener := jsonrpc2test.NetPipeListener()
	cases.Test(t, ctx, listener, jsonrpc2.HeaderFramer(), true)
}

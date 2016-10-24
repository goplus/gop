// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/binary"
	"reflect"
	"testing"
	"time"
)

func TestDecode(t *testing.T) {
	r := new(Recorder)
	stdout := r.Stdout()
	stderr := r.Stderr()

	stdout.Write([]byte("head"))
	stdout.Write(pbWrite(0, "one"))
	stdout.Write(pbWrite(0, "two"))
	stderr.Write(pbWrite(1*time.Second, "three"))
	stdout.Write(pbWrite(2*time.Second, "four"))
	stderr.Write(pbWrite(3*time.Second, "five"))
	stdout.Write([]byte("middle"))
	stderr.Write(pbWrite(4*time.Second, "six"))
	stdout.Write(pbWrite(4*time.Second, "seven"))
	stdout.Write([]byte("tail"))

	want := []Event{
		{"headonetwo", "stdout", 0},
		{"three", "stderr", time.Second},
		{"four", "stdout", time.Second},
		{"five", "stderr", time.Second},
		{"middle", "stdout", 0},
		{"six", "stderr", time.Second},
		{"seventail", "stdout", 0},
	}

	got, err := r.Events()
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func pbWrite(offset time.Duration, s string) []byte {
	out := make([]byte, 16)
	out[2] = 'P'
	out[3] = 'B'
	binary.BigEndian.PutUint64(out[4:], uint64(epoch.Add(offset).UnixNano()))
	binary.BigEndian.PutUint32(out[12:], uint32(len(s)))
	return append(out, s...)
}

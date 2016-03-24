// Copyright 2014 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"
	"unicode/utf8"
)

// When sandbox time begins.
var epoch = time.Unix(1257894000, 0)

// Recorder records the standard and error outputs of a sandbox program
// (comprised of playback headers) and converts it to a sequence of Events.
// It sanitizes each Event's Message to ensure it is valid UTF-8.
//
// Playground programs precede all their writes with a header (described
// below) that describes the time the write occurred (in playground time) and
// the length of the data that will be written. If a non-header is
// encountered where a header is expected, the output is scanned for the next
// header and the intervening text string is added to the sequence an event
// occurring at the same time as the preceding event.
//
// A playback header has this structure:
// 	4 bytes: "\x00\x00PB", a magic header
// 	8 bytes: big-endian int64, unix time in nanoseconds
// 	4 bytes: big-endian int32, length of the next write
//
type Recorder struct {
	mu     sync.Mutex
	writes []recorderWrite
}

type recorderWrite struct {
	b    []byte
	kind string
}

func (r *Recorder) Stdout() io.Writer { return recorderWriter{r, "stdout"} }
func (r *Recorder) Stderr() io.Writer { return recorderWriter{r, "stderr"} }

type recorderWriter struct {
	r    *Recorder
	kind string
}

func (w recorderWriter) Write(b []byte) (n int, err error) {
	w.r.mu.Lock()
	defer w.r.mu.Unlock()

	// Append this write to the previous one if it has the same kind.
	if len(w.r.writes) > 0 {
		prev := &w.r.writes[len(w.r.writes)-1]
		if prev.kind == w.kind {
			prev.b = append(prev.b, b...)
			return len(b), nil
		}
	}

	// Otherwise, append a new write.
	w.r.writes = append(w.r.writes, recorderWrite{
		append([]byte(nil), b...), w.kind,
	})

	return len(b), nil
}

type Event struct {
	Message string        `json:"message"`
	Kind    string        `json:"kind"`  // "stdout" or "stderr"
	Delay   time.Duration `json:"delay"` // time to wait before printing Message
}

func (r *Recorder) Events() ([]Event, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var (
		out []Event
		now = epoch
	)
	for _, w := range r.writes {
		events, err := decode(w.kind, w.b)
		if err != nil {
			return nil, err
		}
		for _, e := range events {
			delay := e.time.Sub(now)
			if delay < 0 {
				delay = 0
			}
			out = append(out, Event{
				Message: string(sanitize(e.msg)),
				Kind:    e.kind,
				Delay:   delay,
			})
			if delay > 0 {
				now = e.time
			}
		}
	}
	return out, nil
}

type event struct {
	msg  []byte
	kind string
	time time.Time
}

func decode(kind string, output []byte) ([]event, error) {
	var (
		magic     = []byte{0, 0, 'P', 'B'}
		headerLen = 8 + 4
		last      = epoch
		events    []event
	)
	add := func(t time.Time, b []byte) {
		var prev *event
		if len(events) > 0 {
			prev = &events[len(events)-1]
		}
		if prev != nil && t.Equal(prev.time) {
			// Merge this event with previous event, to avoid
			// sending a lot of events for a big output with no
			// significant timing information.
			prev.msg = append(prev.msg, b...)
		} else {
			e := event{msg: b, kind: kind, time: t}
			events = append(events, e)
		}
		last = t
	}
	for i := 0; i < len(output); {
		if !bytes.HasPrefix(output[i:], magic) {
			// Not a header; find next header.
			j := bytes.Index(output[i:], magic)
			if j < 0 {
				// No more headers; bail.
				add(last, output[i:])
				break
			}
			add(last, output[i:i+j])
			i += j
		}
		i += len(magic)

		// Decode header.
		if len(output)-i < headerLen {
			return nil, errors.New("short header")
		}
		header := output[i : i+headerLen]
		nanos := int64(binary.BigEndian.Uint64(header[0:]))
		t := time.Unix(0, nanos)
		if t.Before(last) {
			// Force timestamps to be monotonic. (This could
			// be an encoding error, which we ignore now but will
			// will likely be picked up when decoding the length.)
			t = last
		}
		n := int(binary.BigEndian.Uint32(header[8:]))
		if n < 0 {
			return nil, fmt.Errorf("bad length: %v", n)
		}
		i += headerLen

		// Slurp output.
		// Truncated output is OK (probably caused by sandbox limits).
		end := i + n
		if end > len(output) {
			end = len(output)
		}
		add(t, output[i:end])
		i += n
	}
	return events, nil
}

// sanitize scans b for invalid utf8 code points. If found, it reconstructs
// the slice replacing the invalid codes with \uFFFD, properly encoded.
func sanitize(b []byte) []byte {
	if utf8.Valid(b) {
		return b
	}
	var buf bytes.Buffer
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		b = b[size:]
		buf.WriteRune(r)
	}
	return buf.Bytes()
}

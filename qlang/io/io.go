package io

import (
	"io"
)

// -----------------------------------------------------------------------------

// Exports is the export table of this module.
//
var Exports = map[string]interface{}{
	"_name":       "io",
	"copy":        io.Copy,
	"copyN":       io.CopyN,
	"readAtLeast": io.ReadAtLeast,
	"readFull":    io.ReadFull,
	"writeString": io.WriteString,

	"Copy":        io.Copy,
	"CopyN":       io.CopyN,
	"ReadAtLeast": io.ReadAtLeast,
	"ReadFull":    io.ReadFull,
	"WriteString": io.WriteString,

	"pipe":          io.Pipe,
	"limitReader":   io.LimitReader,
	"multiReader":   io.MultiReader,
	"multiWriter":   io.MultiWriter,
	"teeReader":     io.TeeReader,
	"sectionReader": io.NewSectionReader,

	"Pipe":             io.Pipe,
	"LimitReader":      io.LimitReader,
	"MultiReader":      io.MultiReader,
	"MultiWriter":      io.MultiWriter,
	"TeeReader":        io.TeeReader,
	"NewSectionReader": io.NewSectionReader,

	"EOF":              io.EOF,
	"ErrClosedPipe":    io.ErrClosedPipe,
	"ErrNoProgress":    io.ErrNoProgress,
	"ErrShortBuffer":   io.ErrShortBuffer,
	"ErrShortWrite":    io.ErrShortWrite,
	"ErrUnexpectedEOF": io.ErrUnexpectedEOF,
}

// -----------------------------------------------------------------------------

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package base defines shared basic pieces of the gop command,
// in particular logging and the Command structure.
package base

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// A Command is an implementation of a gop command
// like gop export or gop install.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The words between "gop" and the first flag or argument in the line are taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'gop help' output.
	Short string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// Commands lists the available commands and help topics.
	// The order here is the order in which they are printed by 'gop help'.
	// Note that subcommands are in general best avoided.
	Commands []*Command
}

// Gop command
var Gop = &Command{
	UsageLine: "gop",
	Short:     `Gop is a tool for managing Go+ source code.`,
	// Commands initialized in package main
}

// LongName returns the command's long name: all the words in the usage line between "gop" and a flag or argument,
func (c *Command) LongName() string {
	name := c.UsageLine
	if i := strings.Index(name, " ["); i >= 0 {
		name = name[:i]
	}
	if name == "gop" {
		return ""
	}
	return strings.TrimPrefix(name, "gop ")
}

// Name returns the command's short name: the last word in the usage line before a flag or argument.
func (c *Command) Name() string {
	name := c.LongName()
	if i := strings.LastIndex(name, " "); i >= 0 {
		name = name[i+1:]
	}
	return name
}

// Usage show the command usage.
func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s\n", c.Short, c.UsageLine)
	c.Flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// Usage is the usage-reporting function, filled in by package main
// but here for reference by other packages.
var Usage func()

// CmdName - "build", "install", "list", "mod tidy", etc.
var CmdName string

// Main runs a command.
func Main(c *Command, app string, args []string) {
	name := c.UsageLine
	if i := strings.Index(name, " ["); i >= 0 {
		c.UsageLine = app + name[i:]
	}
	c.Run(c, args)
}

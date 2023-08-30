package test

import (
	"github.com/goplus/gop/cmd/internal/base"
)

type boolFlag interface {
	IsBoolFlag() bool
}

func PassTestFlags(cmd *base.Command) *base.PassArgs {
	p := base.PassBuildFlags(cmd)
	p.Bool("c", "i", "cover", "json", "benchmem", "failfast", "short")
	p.Var("o", "covermode", "coverpkg", "exec", "vet",
		"bench", "benchtime", "blockprofile", "blockprofilerate",
		"count", "coverprofile", "cpu", "cpuprofile",
		"fuzz", "list", "memprofile", "memprofilerate",
		"mutexprofile", "mutexprofilefraction", "outputdir", "parallel",
		"run", "timeout", "fuzztime", "fuzzminimizetime",
		"trace", "shuffle")
	for name := range passFlagToTest {
		if b, ok := cmd.Flag.Lookup(name).Value.(boolFlag); ok && b.IsBoolFlag() {
			p.Bool("test." + name)
		} else {
			p.Var("test." + name)
		}
	}
	return p
}

var passFlagToTest = map[string]bool{
	"bench":                true,
	"benchmem":             true,
	"benchtime":            true,
	"blockprofile":         true,
	"blockprofilerate":     true,
	"count":                true,
	"coverprofile":         true,
	"cpu":                  true,
	"cpuprofile":           true,
	"failfast":             true,
	"fuzz":                 true,
	"fuzzminimizetime":     true,
	"fuzztime":             true,
	"list":                 true,
	"memprofile":           true,
	"memprofilerate":       true,
	"mutexprofile":         true,
	"mutexprofilefraction": true,
	"outputdir":            true,
	"parallel":             true,
	"run":                  true,
	"short":                true,
	"shuffle":              true,
	"timeout":              true,
	"trace":                true,
	"v":                    true,
}

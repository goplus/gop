package modfile

import (
	"context"
	"testing"

	"golang.org/x/mod/modfile"
)

func fixVersion(ctx context.Context, fixed *bool) modfile.VersionFixer {
	return func(path, vers string) (resolved string, err error) {
		return vers, nil
	}
}

var gopmod = `
module spx

go 1.17
gop 1.1

classfile .gmx .spx github.com/goplus/spx

require (
    github.com/ajstarks/svgo v0.0.0-20210927141636-6d70534b1098
)
`

var gopmod2 = `
module moduleUserProj

go 1.17
gop 1.1

register spx

require (
    github.com/goplus/spx v1.0
)
`

func TestParse(t *testing.T) {
	f, err := Parse("github.com/goplus/gop/gop.mod", []byte(gopmod), func(path, vers string) (resolved string, err error) {
		return vers, nil
	})
	if err != nil {
		t.Error(err)
	}
	if f.Gop.Version != "1.1" {
		t.Errorf("gop version expected be 1.1, but %s got", f.Gop.Version)
	}

	if len(f.Classfile.Exts) != 2 {
		t.Errorf("classfile exts length expected be 2, but %d got", len(f.Classfile.Exts))
	}

	if f.Classfile.Exts[0] != ".gmx" {
		t.Errorf("classfile exts expected be .gmx, but %s got", f.Classfile.Exts[0])
	}
	if f.Classfile.Exts[1] != ".spx" {
		t.Errorf("classfile exts expected be .spx, but %s got", f.Classfile.Exts[0])
	}

	if f.Classfile.Path != "github.com/goplus/spx" {
		t.Errorf("classfile path expected be github.com/goplus/spx, but %s got", f.Classfile.Path)
	}

	f2, err := Parse("github.com/goplus/gop/gop.mod", []byte(gopmod2), func(path, vers string) (resolved string, err error) {
		return vers, nil
	})

	if err != nil {
		t.Error(err)
	}

	if f2.Register.ClassfileMod != "spx" {
		t.Errorf("register classfile mod expected be spx, but %s got", f2.Register.ClassfileMod)
	}
}

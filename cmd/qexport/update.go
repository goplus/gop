package main

import (
	"bytes"
	"go/ast"
	"go/build"
	"go/format"
	"go/token"
	"go/types"
	"io/ioutil"
	"path/filepath"

	"github.com/visualfc/gotools/pkgwalk"
)

type ObjectInfo struct {
	kind  pkgwalk.ObjKind
	ident *ast.Ident
	obj   types.Object
}

type InsertInfo struct {
	file           *ast.File
	decl           ast.Decl
	pos            token.Position
	specImportName string
	importPos      token.Position
}

type UpdateInfo struct {
	fileset     *token.FileSet
	updataPkg   string
	usesMap     map[string]*ObjectInfo
	exportsInfo *InsertInfo
	initsInfo   []*InsertInfo
}

func insertData(org []byte, data []byte, offset int) []byte {
	var out []byte
	out = append(out, org[:offset-1]...)
	out = append(out, data...)
	out = append(out, org[offset-1:]...)
	return out
}

func (i *InsertInfo) UpdateFile(outpath string, data []byte, hasTypeExport bool) error {
	file := filepath.Join(outpath, i.pos.Filename)
	all, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if hasTypeExport && i.specImportName == "spec" {
		data = bytes.Replace(data, []byte("qlang."), []byte("spec."), -1)
	}

	out := insertData(all, data, i.pos.Offset)
	if hasTypeExport && i.specImportName == "" {
		spec := []byte("\n\nqlang \"qlang.io/spec\"\n")
		out = insertData(out, spec, i.importPos.Offset+1)
	}

	// format
	fout, err := format.Source(out)
	if err != nil {
		return err
	}

	ioutil.WriteFile(file, fout, 0777)
	return nil
}

//check has update package
func CheckUpdateInfo(pkgname string, pkg string) (*UpdateInfo, error) {
	updatePkgPath := flagUpdatePath + "/" + pkg
	bp, err := build.Import(updatePkgPath, "", 0)
	if err != nil {
		return nil, err
	}
	CopyDir(bp.Dir, flagExportPath+"/"+pkg, false)

	conf := pkgwalk.DefaultPkgConfig()
	w := pkgwalk.NewPkgWalker(&build.Default)
	pkgx, err := pkgwalk.ImportPackage(w, updatePkgPath, conf)
	if err != nil {
		return nil, err
	}
	list := pkgwalk.LookupObjList(w, pkgx, conf)

	ui := &UpdateInfo{}
	ui.fileset = w.FileSet
	ui.updataPkg = bp.ImportPath
	ui.usesMap = make(map[string]*ObjectInfo)

	for ident, obj := range conf.Info.Uses {
		if obj != nil && obj.Pkg() != nil && obj.Pkg().Path() == pkg {
			kind, _ := pkgwalk.ParserObjectKind(ident, obj, conf)
			ui.usesMap[ident.Name] = &ObjectInfo{kind, ident, obj}
		}
	}
	fnMakeInsertInfo := func(ident *ast.Ident) *InsertInfo {
		file, decl := w.FindDeclForPos(ident.Pos())
		if decl == nil {
			return nil
		}
		specImport := w.FindImportName(file, "qlang.io/spec")
		endImportPos := w.FindImportEndPos(file)
		pos := w.FileSet.Position(decl.End())
		_, pos.Filename = filepath.Split(pos.Filename)
		ipos := w.FileSet.Position(endImportPos)
		_, ipos.Filename = filepath.Split(pos.Filename)
		return &InsertInfo{file, decl, pos, specImport, ipos}
	}
	for _, obj := range list {
		if obj != nil && obj.Obj != nil {
			if obj.Obj.Name() == "Exports" {
				v := fnMakeInsertInfo(obj.Ident)
				if v != nil {
					ui.exportsInfo = v
				}
			} else if obj.Obj.Name() == "init" {
				v := fnMakeInsertInfo(obj.Ident)
				if v != nil {
					ui.initsInfo = append(ui.initsInfo, v)
				}
			}
		}
	}

	return ui, nil
}

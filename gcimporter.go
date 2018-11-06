// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gcimporter implements Import for gc-generated object files.
package gcimporter

import (
	"bufio"
	"fmt"
	"go/build"
	"go/token"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// debugging/development support
const debug = false

var pkgExts = [...]string{".a", ".o"}

func findPkgContext(
	ctx *build.Context, path, srcDir string,
) (filename, id string) {
	if path == "" {
		return
	}

	if build.IsLocalImport(path) || filepath.IsAbs(path) {
		return
	}

	// "x" -> "$GOPATH/pkg/$GOOS_$GOARCH/x.ext", "x"
	// Don't require the source files to be present.
	if abs, err := filepath.Abs(srcDir); err == nil { // see issue 14282
		srcDir = abs
	}
	bp, _ := ctx.Import(path, srcDir, build.FindOnly|build.AllowBinary)
	if bp.PkgObj == "" {
		id = path // make sure we have an id to print in error message
		return
	}
	noext := strings.TrimSuffix(bp.PkgObj, ".a")
	id = bp.ImportPath

	// try extensions
	for _, ext := range pkgExts {
		filename = noext + ext
		if f, err := os.Stat(filename); err == nil && !f.IsDir() {
			return
		}
	}

	filename = "" // not found
	return
}

// Import imports a gc-generated package given its import path and srcDir, adds
// the corresponding package object to the packages map, and returns the object.
// The packages map must contain all packages already imported.
//
func Import(
	packages map[string]*types.Package, path, srcDir string,
) (pkg *types.Package, err error) {
	return importContext(&build.Default, packages, path, srcDir)
}

func importContext(
	ctx *build.Context,
	packages map[string]*types.Package, path, srcDir string,
) (pkg *types.Package, err error) {
	filename, id := findPkgContext(ctx, path, srcDir)
	if filename == "" {
		if path == "unsafe" {
			return types.Unsafe, nil
		}
		return nil, fmt.Errorf("can't find import: %q", id)
	}

	// no need to re-import if the package was imported completely before
	if pkg = packages[id]; pkg != nil && pkg.Complete() {
		return
	}

	// open file
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			// add file name to error
			err = fmt.Errorf("%s: %v", filename, err)
		}
	}()
	defer f.Close()

	var hdr string
	buf := bufio.NewReader(f)
	if hdr, err = FindExportData(buf); err != nil {
		return
	}

	switch hdr {
	case "$$\n":
		err = fmt.Errorf("import %q: old export format no longer supported (recompile library)", path)

	case "$$B\n":
		var data []byte
		data, err = ioutil.ReadAll(buf)
		if err != nil {
			break
		}

		// TODO(gri): allow clients of go/importer to provide a FileSet.
		// Or, define a new standard go/types/gcexportdata package.
		fset := token.NewFileSet()

		// The indexed export format starts with an 'i'; the older
		// binary export format starts with a 'c', 'd', or 'v'
		// (from "version"). Select appropriate importer.
		if len(data) > 0 && data[0] == 'i' {
			_, pkg, err = iImportData(fset, packages, data[1:], id)
		} else {
			_, pkg, err = BImportData(fset, packages, data, id)
		}

	default:
		err = fmt.Errorf("unknown export data header: %q", hdr)
	}

	return
}

func deref(typ types.Type) types.Type {
	if p, _ := typ.(*types.Pointer); p != nil {
		return p.Elem()
	}
	return typ
}

type byPath []*types.Package

func (a byPath) Len() int           { return len(a) }
func (a byPath) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPath) Less(i, j int) bool { return a[i].Path() < a[j].Path() }

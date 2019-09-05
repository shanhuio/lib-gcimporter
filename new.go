package gcimporter

import (
	"go/build"
	"go/token"
	"go/types"
	"log"
)

// Importer is a GC importer that supports a provided context.
type Importer struct {
	ctx      *build.Context
	packages map[string]*types.Package
	alias    *AliasMap
}

// New makes a new Importer using the given build context.
// It implements go/types.Importer
func New(ctx *build.Context, alias *AliasMap) *Importer {
	return &Importer{
		ctx:      ctx,
		packages: make(map[string]*types.Package),
		alias:    alias,
	}
}

func (imp *Importer) mapPath(path string) string {
	if imp.alias != nil {
		return imp.alias.Map(path)
	}
	return path
}

// Import imports a given package of the path.
func (imp *Importer) Import(path string) (*types.Package, error) {
	p, err := imp.ImportFrom(path, "", 0)
	if err != nil {
		log.Printf("import %q: %s", path, err)
	}
	return p, err
}

// ImportFrom imports a given package of the path at the source directory.
// mode must be 0.
func (imp *Importer) ImportFrom(path, srcDir string, mode types.ImportMode) (
	*types.Package, error,
) {
	if mode != 0 {
		panic("mode must be 0")
	}
	mapped := imp.mapPath(path)
	fset := token.NewFileSet()
	p, err := importContext(imp.ctx, fset, imp.packages, mapped, srcDir, nil)
	if err != nil {
		log.Printf("importFrom %q(%q), %q: %s", path, mapped, srcDir, err)
	}
	return p, err
}

package gcimporter

import (
	"go/build"
	goimp "go/importer"
	"go/types"
	"log"
)

// Importer is a GC importer that supports a provided context.
type Importer struct {
	ctx      *build.Context
	packages map[string]*types.Package
	lookup   goimp.Lookup
}

// New makes a new Importer using the given build context.
// It implements go/types.Importer
func New(ctx *build.Context, lookup goimp.Lookup) *Importer {
	return &Importer{
		ctx:      ctx,
		packages: make(map[string]*types.Package),
		lookup:   lookup,
	}
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
	p, err := importContext(imp.ctx, imp.packages, path, srcDir, imp.lookup)
	if err != nil {
		log.Printf("importFrom %q, %q: %s", path, srcDir, err)
	}
	return p, err
}

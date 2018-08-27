package gcimporter

import (
	"go/build"
	"go/types"
)

// Importer is a GC importer that supports a provided context.
type Importer struct {
	ctx      *build.Context
	packages map[string]*types.Package
}

// New makes a new Importer using the given build context.
// It implements go/types.Importer
func New(ctx *build.Context) *Importer {
	return &Importer{
		ctx:      ctx,
		packages: make(map[string]*types.Package),
	}
}

// Import imports a given package of the path.
func (imp *Importer) Import(path string) (*types.Package, error) {
	return imp.ImportFrom(path, "", 0)
}

// ImportFrom imports a given package of the path at the source directory.
// mode must be 0.
func (imp *Importer) ImportFrom(path, srcDir string, mode types.ImportMode) (
	*types.Package, error,
) {
	if mode != 0 {
		panic("mode must be 0")
	}
	return importContext(imp.ctx, imp.packages, path, srcDir, nil)
}

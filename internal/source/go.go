package source

import (
	"go/build"
)

type GoPackage struct {
	Dir      string
	HasTests bool
}

func FindGoPackage(rootPath string, dir string) (GoPackage, error) {
	ctx := build.Default
	p, err := ctx.Import(dir, rootPath, 0)
	if err != nil {
		return GoPackage{}, err
	}
	testFiles := append(p.TestGoFiles, p.XTestGoFiles...)
	return GoPackage{
		Dir:      p.Dir,
		HasTests: len(testFiles) > 0,
	}, nil
}

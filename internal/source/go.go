package source

import (
	"go/build"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

type GoPackage struct {
	Dir      string
	HasTests bool
}

type Finder struct {
	bctx build.Context
}

func NewFinder() *Finder {
	bctx := build.Default
	return &Finder{
		bctx: bctx,
	}
}

func (f *Finder) Package(rootPath string, dir string) (GoPackage, error) {
	p, err := f.bctx.Import(dir, rootPath, 0)
	if err != nil {
		return GoPackage{}, err
	}
	testFiles := append(p.TestGoFiles, p.XTestGoFiles...)
	return GoPackage{
		Dir:      p.Dir,
		HasTests: len(testFiles) > 0,
	}, nil
}

func FindModuleDirs(rootPath string) map[string]string {
	modules := map[string]string{}
	filepath.Walk(rootPath, func(file string, info fs.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		// Skip vendor directory
		if info.IsDir() && info.Name() == "vendor" {
			os.Stat(path.Join(info.Name(), "modules.txt"))
			return filepath.SkipDir
		}

		if !info.IsDir() && info.Name() == "go.mod" {
			b, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			p := modfile.ModulePath(b)
			if p != "" {
				modules[p] = filepath.Dir(file)
			}
		}
		return nil
	})
	return modules
}

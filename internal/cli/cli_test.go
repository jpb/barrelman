package cli_test

import (
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jpb/barrelman/internal/cli"
	"github.com/stretchr/testify/require"
)

//go:embed test/*
var f embed.FS

func read(t *testing.T, path string) io.Reader {
	t.Helper()
	b, err := f.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(b)
}

type rewrittenFile struct {
	fs.FileInfo
}

func (f rewrittenFile) Name() string {
	return strings.TrimPrefix(f.FileInfo.Name(), "_")
}

func readDir(path string) ([]fs.FileInfo, error) {
	dir, err := f.ReadDir(path)
	infos := []fs.FileInfo{}
	if err != nil {
		return nil, err
	}
	for _, d := range dir {
		info, err := d.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, rewrittenFile{info})
		fmt.Println(rewrittenFile{info}.Name())
	}
	return infos, nil
}

func openFile(filepath string) (io.ReadCloser, error) {
	filename := path.Base(filepath)
	if filename == "go.mod" || filename == "go.sum" {
		dir := path.Dir(filepath)
		filepath = path.Join(dir, "_"+filename)
	}
	return f.Open(filepath)
}

func isDir(dir string) bool {
	f, err := f.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.IsDir()
}

func hasSubdir(root, dir string) (rel string, ok bool) {
	const sep = string(filepath.Separator)
	root = filepath.Clean(root)
	if !strings.HasSuffix(root, sep) {
		root += sep
	}
	dir = filepath.Clean(dir)
	after, found := strings.CutPrefix(dir, root)
	if !found {
		return "", false
	}
	return filepath.ToSlash(after), true
}

func rootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "../..")
}

func TestRun(t *testing.T) {
	require := require.New(t)
	actual := []cli.Warning{}
	cli.Run(
		rootDir(),
		read(t, "test/coverage.out"),
		read(t, "test/additions.diff"),
		func(w cli.Warning) {
			actual = append(actual, w)
		},
	)
	require.Equal([]cli.Warning{
		{
			Filepath:  "test/adder/adder.go",
			StartLine: 10,
			EndLine:   11,
			Message:   "Lines 10-11 have no test coverage.",
		},
	}, actual)
}

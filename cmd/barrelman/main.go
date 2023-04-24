package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jpb/barrelman/internal/cli"
	"github.com/jpb/barrelman/internal/report"
)

func main() {
	root, ok := os.LookupEnv("BARRELMAN_ROOT_DIR")
	if !ok {
		error("BARRELMAN_ROOT_DIR must be set")
	}
	diffPath, ok := os.LookupEnv("BARRELMAN_DIFF")
	if !ok {
		error("BARRELMAN_DIFF must be set")
	}
	coverPath, ok := os.LookupEnv("BARRELMAN_COVER")
	if !ok {
		error("BARRELMAN_COVER must be set")
	}

	root, err := filepath.Abs(root)
	if err != nil {
		error(err)
	}
	err = cli.Run(
		root,
		read(coverPath),
		read(diffPath),
		report.GitHubActionAnnotation,
	)
	if err != nil {
		error(err)
	}
}

func read(path string) io.Reader {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func error(message interface{}) {
	fmt.Println(message)
	os.Exit(1)
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jpb/barrelman/internal/cli"
	"github.com/jpb/barrelman/internal/report"
)

func main() {
	root, ok := os.LookupEnv("BARRELMAN_ROOT_DIR")
	if !ok {
		var err error
		root, err = gitRoot()
		if err != nil {
			fatal(
				fmt.Sprintf("unable to find git root: %v, please provide BARRELMAN_ROOT_DIR", err),
			)
		}
	}
	if len(os.Args) != 3 {
		fatal("usage: %s coverage.out <(git diff ...)", os.Args[0])
	}
	coverPath := os.Args[1]
	diffPath := os.Args[2]

	root, err := filepath.Abs(root)
	if err != nil {
		fatal(err.Error())
	}
	err = cli.Run(
		root,
		read(coverPath),
		read(diffPath),
		report.GitHubActionAnnotation,
	)
	if err != nil {
		fatal(err.Error())
	}
}

func read(path string) io.Reader {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(b)
}

func fatal(message string, args ...interface{}) {
	fmt.Printf(message, args...)
	fmt.Println()
	os.Exit(1)
}

func gitRoot() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}

package cli

import (
	"io"
	"path"
	"path/filepath"

	"github.com/jpb/barrelman/internal/cover"
	"github.com/jpb/barrelman/internal/git"
	"github.com/jpb/barrelman/internal/hunk"
	"github.com/jpb/barrelman/internal/source"
)

type Warning struct {
	Filepath  string
	StartLine int
	EndLine   int
	Message   string
}

// Run runs the CLI
func Run(
	rootPath string,
	coverProfile io.Reader,
	diff io.Reader,
	report func(Warning),
) error {
	additions, err := git.Additions(diff)
	if err != nil {
		return err
	}
	// Make paths absolute
	for i, hunk := range additions {
		additions[i].Filepath = path.Join(rootPath, hunk.Filepath)
	}
	additionsByFile := fileHunks(additions)

	uncovered, err := cover.Uncovered(coverProfile)
	if err != nil {
		return err
	}
	// The "file" in a coverage profile is actually the go import name, and a go
	// import can include a module name. Rewrite the path to be relative to the
	// root of the project (as git.Additions does)
	for i, hunk := range uncovered {
		dir := path.Dir(hunk.Filepath)
		p, err := source.FindGoPackage(rootPath, dir)
		uncovered[i].Filepath = path.Join(p.Dir, filepath.Base(hunk.Filepath))
		if err != nil {
			return err
		}
	}
	uncoveredByFile := fileHunks(uncovered)

	for file, uncovered := range uncoveredByFile {
		additions, ok := additionsByFile[file]
		if !ok {
			// The uncovered file isn't in the additions
			continue
		}

		// Find intersections between additions
		for _, hunk := range hunk.Intersections(additions, uncovered) {
			path, _ := filepath.Rel(rootPath, file)
			report(Warning{
				Filepath:  path,
				StartLine: hunk.StartLine,
				EndLine:   hunk.EndLine,
				Message:   "no test coverage",
			})
		}
	}

	// Find files with no test coverage at all
	for _, hunks := range additionsByFile {
		hunk := hunks[0]
		relative, _ := filepath.Rel(rootPath, hunk.Filepath)
		dir := "./" + path.Dir(relative)
		p, err := source.FindGoPackage(rootPath, dir)
		if err != nil {
			return err
		}
		if !p.HasTests {
			report(Warning{
				Filepath:  hunk.Filepath,
				StartLine: hunk.StartLine,
				Message:   "no test files found for this package",
			})
		}
	}

	return nil
}

func fileHunks(hunks []hunk.Hunk) map[string][]hunk.Hunk {
	m := map[string][]hunk.Hunk{}
	hunk.SortHunks(hunks)
	for _, hunk := range hunks {
		m[hunk.Filepath] = append(m[hunk.Filepath], hunk)
	}
	return m
}

package cli

import (
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"

	"github.com/jpb/barrelman/internal/cover"
	"github.com/jpb/barrelman/internal/git"
	"github.com/jpb/barrelman/internal/hunk"
	"github.com/jpb/barrelman/internal/source"
)

// Warning represents a issue that barrel man has discovered. For human
// consumption.
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
	// All path comparisons are absolute

	finder := source.NewFinder()
	moduleDirs := source.FindModuleDirs(rootPath)

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
		// Replace the module name in the path with the directory location of the module
		for module, dir := range moduleDirs {
			if strings.HasPrefix(hunk.Filepath, module) {
				uncovered[i].Filepath = strings.Replace(
					hunk.Filepath,
					module,
					dir,
					1,
				)
				break
			}
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
			msg := fmt.Sprintf("Line %d has no test coverage.", hunk.StartLine)
			if hunk.StartLine != hunk.EndLine {
				msg = fmt.Sprintf("Lines %d-%d have no test coverage.", hunk.StartLine, hunk.EndLine)
			}
			report(Warning{
				Filepath:  path,
				StartLine: hunk.StartLine,
				EndLine:   hunk.EndLine,
				Message:   msg,
			})
		}
	}

	// Find files with no test coverage at all
	for file, hunks := range additionsByFile {
		if _, ok := uncoveredByFile[file]; ok {
			// Skip this file as it would be handled above
			continue
		}
		hunk := hunks[0]
		relative, _ := filepath.Rel(rootPath, hunk.Filepath)
		dir := "./" + path.Dir(relative)
		p, err := finder.Package(rootPath, dir)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !p.HasTests {
			report(Warning{
				Filepath:  hunk.Filepath,
				StartLine: hunk.StartLine,
				Message:   "There were no test files found for this package.",
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

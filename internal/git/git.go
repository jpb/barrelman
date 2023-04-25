package git

import (
	"io"

	"github.com/bluekeyes/go-gitdiff/gitdiff"

	"github.com/jpb/barrelman/internal/hunk"
)

// Additions will return the hunks that were added in the provided diff
func Additions(diff io.Reader) ([]hunk.Hunk, error) {
	files, _, err := gitdiff.Parse(diff)
	if err != nil {
		return nil, err
	}

	hunks := []hunk.Hunk{}

	for _, file := range files {
		for _, t := range file.TextFragments {
			if t.LinesAdded == 0 {
				continue
			}
			hunks = append(hunks, hunk.Hunk{
				Filepath:  file.NewName,
				StartLine: int(t.NewPosition),
				EndLine:   int(t.NewPosition + t.LinesAdded - 1),
			})
		}
	}

	return hunks, nil
}

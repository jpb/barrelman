package git

import (
	"io"

	"github.com/bluekeyes/go-gitdiff/gitdiff"

	"github.com/jpb/barrelman/internal/data"
)

// Additions will return the hunks that were added in the provided diff
func Additions(diff io.Reader) ([]data.Hunk, error) {
	files, _, err := gitdiff.Parse(diff)
	if err != nil {
		return nil, err
	}

	hunks := []data.Hunk{}

	for _, file := range files {
		for _, t := range file.TextFragments {
			if t.LinesAdded == 0 {
				continue
			}
			hunks = append(hunks, data.Hunk{
				Filepath:  file.NewName,
				StartLine: int(t.NewPosition),
				EndLine:   int(t.NewPosition + t.LinesAdded - 1),
			})
		}
	}

	return hunks, nil
}

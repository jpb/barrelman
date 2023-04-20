package cover

import (
	"io"

	"golang.org/x/tools/cover"

	"github.com/jpb/barrelman/internal/hunk"
)

// Uncovered returns the hunks with no coverage
func Uncovered(r io.Reader) ([]hunk.Hunk, error) {
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		return nil, err
	}

	hunks := []hunk.Hunk{}

	for _, profile := range profiles {
		for _, block := range profile.Blocks {
			if block.Count > 0 {
				continue
			}
			hunks = append(hunks, hunk.Hunk{
				Filepath:  profile.FileName,
				StartLine: block.StartLine,
				EndLine:   block.EndLine,
			})
		}
	}

	return hunk.CompressHunks(hunks), nil
}

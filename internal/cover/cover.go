package cover

import (
	"io"
	"strings"

	"golang.org/x/tools/cover"

	"github.com/jpb/barrelman/internal/data"
)

// Uncovered returns the hunks with no coverage
func Uncovered(packageName string, r io.Reader) ([]data.Hunk, error) {
	profiles, err := cover.ParseProfilesFromReader(r)
	if err != nil {
		return nil, err
	}

	hunks := []data.Hunk{}

	for _, profile := range profiles {
		for _, block := range profile.Blocks {
			if block.Count > 0 {
				continue
			}
			hunks = append(hunks, data.Hunk{
				Filepath: strings.TrimPrefix(
					strings.TrimPrefix(profile.FileName, packageName),
					"/",
				),
				StartLine: block.StartLine,
				EndLine:   block.EndLine,
			})
		}
	}

	return data.CompressHunks(hunks), nil
}

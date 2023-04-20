package data

import "sort"

// Hunk describes a section of source code
type Hunk struct {
	Filepath  string
	StartLine int
	EndLine   int
}

// CompressHunks joins hunks that adjacent in the same file
func CompressHunks(hunks []Hunk) []Hunk {
	// Copy so sort does not modify hunks
	in := make([]Hunk, len(hunks))
	copy(in, hunks)

	// Pull all hunks in order
	sort.Slice(in, func(i, j int) bool {
		ii := in[i]
		jj := in[j]
		if ii.Filepath == jj.Filepath {
			return ii.StartLine < jj.StartLine
		}
		return ii.Filepath < jj.Filepath
	})

	out := []Hunk{}
	prev := Hunk{}

	for i := 0; i < len(in); i++ {
		// If file and start/end lines match, just extending the last hunk
		if in[i].Filepath == prev.Filepath {
			if in[i].StartLine == prev.EndLine || in[i].StartLine == prev.EndLine+1 {
				prev.EndLine = in[i].EndLine
				out[len(out)-1] = prev
				continue
			}
		}

		// Otherwise keep it as-is
		out = append(out, in[i])
		prev = in[i]
	}

	return out
}

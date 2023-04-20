package hunk

import (
	"sort"
)

// Hunk describes a section of source code
type Hunk struct {
	Filepath  string
	StartLine int
	EndLine   int
}

// CompressHunks joins hunks that adjacent in the same file
func CompressHunks(hs []Hunk) []Hunk {
	// Copy so sort does not modify hunks
	in := make([]Hunk, len(hs))
	copy(in, hs)

	// Pull all hunks in order
	SortHunks(in)

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

// SortHunks sorts hunks by file and location in file
func SortHunks(hunks []Hunk) {
	sort.Slice(hunks, func(i, j int) bool {
		ii := hunks[i]
		jj := hunks[j]
		if ii.Filepath == jj.Filepath {
			return ii.StartLine < jj.StartLine
		}
		return ii.Filepath < jj.Filepath
	})
}

// Intersections returns a list of hunks of lines that are present in both as
// and bs.
// Assumptions:
//   - as and bs are sorted
//   - as and bs do not themselves have any overlapping hunks
//   - EndLine >= StartLine for each hunk
//   - Filepath is the same for all hunks
func Intersections(as, bs []Hunk) []Hunk {
	hunks := []Hunk{}
	aIdx := 0
	bIdx := 0

	for aIdx < len(as) && bIdx < len(bs) {
		a := as[aIdx]
		b := bs[bIdx]

		// Do they overlap?
		// if a start is within start and end
		if (a.StartLine >= b.StartLine && a.StartLine <= b.EndLine) ||
			(b.StartLine >= a.StartLine && b.StartLine <= a.EndLine) {
			hunks = append(hunks, Hunk{
				Filepath:  a.Filepath,
				StartLine: max(a.StartLine, b.StartLine),
				EndLine:   min(a.EndLine, b.EndLine),
			})
		}

		// advance the smaller one
		if a.EndLine < b.EndLine {
			aIdx++
		} else {
			bIdx++
		}
	}

	return hunks
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

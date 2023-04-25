package hunk_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jpb/barrelman/internal/hunk"
)

func TestIntersections(t *testing.T) {
	require := require.New(t)
	as := []hunk.Hunk{
		{StartLine: 1, EndLine: 2},
		{StartLine: 5, EndLine: 9},
		{StartLine: 9, EndLine: 11},
		{StartLine: 100, EndLine: 101},
		{StartLine: 190, EndLine: 210},
		{StartLine: 300, EndLine: 305},
		{StartLine: 405, EndLine: 409},
	}
	bs := []hunk.Hunk{
		{StartLine: 1, EndLine: 2},
		{StartLine: 4, EndLine: 6},
		{StartLine: 90, EndLine: 110},
		{StartLine: 200, EndLine: 201},
		{StartLine: 305, EndLine: 309},
		{StartLine: 400, EndLine: 405},
		{StartLine: 10000, EndLine: 10000},
	}
	expected := []hunk.Hunk{
		{StartLine: 1, EndLine: 2},
		{StartLine: 5, EndLine: 6},
		{StartLine: 100, EndLine: 101},
		{StartLine: 200, EndLine: 201},
		{StartLine: 305, EndLine: 305},
		{StartLine: 405, EndLine: 405},
	}

	actual := hunk.Intersections(as, bs)

	require.Equal(expected, actual)
}

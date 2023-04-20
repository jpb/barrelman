package cover_test

import (
	"embed"
	"testing"

	"github.com/jpb/barrelman/internal/cover"
	"github.com/jpb/barrelman/internal/hunk"
	"github.com/stretchr/testify/require"
)

//go:embed test/*
var f embed.FS

func TestUncovered(t *testing.T) {
	require := require.New(t)
	file, err := f.Open("test/coverage.out")
	require.Nil(err)
	defer file.Close()
	hunks, err := cover.Uncovered(file)

	require.Equal(
		[]hunk.Hunk{
			{Filepath: "github.com/jpb/barrelman/test/fizzbuzz/fizzbuzz.go", StartLine: 15, EndLine: 17},
			{Filepath: "github.com/jpb/barrelman/test/fizzbuzz/fizzbuzz.go", StartLine: 19, EndLine: 22},
		},
		hunks,
	)
}

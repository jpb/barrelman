package git_test

import (
	"embed"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jpb/barrelman/internal/git"
)

//go:embed test/*.diff
var f embed.FS

func TestAdditions(t *testing.T) {
	require := require.New(t)

	noAdditions := map[string]struct{}{
		"removefile.diff": struct{}{},
		"deletion.diff":   struct{}{},
	}

	files, err := f.ReadDir("test")
	require.Nil(err)
	for _, file := range files {
		a, err := f.Open(path.Join("test", file.Name()))
		require.Nil(err)
		hunks, err := git.Additions(a)
		a.Close()
		require.Nil(err)

		_, ok := noAdditions[file.Name()]
		if ok {
			require.Len(hunks, 0)
		} else {
			require.Greater(len(hunks), 0)
		}

		for _, hunk := range hunks {
			require.True(hunk.StartLine <= hunk.EndLine)
			require.True(
				strings.HasPrefix(hunk.Filepath, "test/fizzbuzz/fizzbuzz"),
			)
			require.True(
				strings.HasSuffix(hunk.Filepath, ".go"),
			)
		}
	}
}

package source_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jpb/barrelman/internal/source"
)

func TestFindModules(t *testing.T) {
	require := require.New(t)
	rootDir, err := filepath.Abs("../..")
	require.Nil(err)
	expected := map[string]string{
		"github.com/jpb/barrelman": rootDir,
	}
	actual := source.FindModuleDirs(rootDir)
	require.Equal(expected, actual)
}

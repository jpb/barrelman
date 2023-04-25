package adder_test

import (
	"testing"

	"github.com/jpb/barrelman/test/adder"
	"github.com/stretchr/testify/require"
)

func TestAdder(t *testing.T) {
	require := require.New(t)
	require.Equal(6, adder.AddNums(1, 2, 3))
}

package fizzbuzz_test

import (
	"testing"

	"github.com/jpb/barrelman/test/fizzbuzz"
	"github.com/stretchr/testify/require"
)

func TestFizzBuzz(t *testing.T) {
	require := require.New(t)
	require.Equal("fizz", fizzbuzz.FizzBuzz(3))
}

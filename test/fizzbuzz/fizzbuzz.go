package fizzbuzz

import (
	"strconv"
)

const (
	fizz = "fizz"
	buzz = "buzz"
)

// FizzBuzz returns "fizz" if i is divisible by 3, "buzz" if i is divisible by
// 5, "fizzbuzz" if i is divisible by both 3 and 5, and i as a string otherwise
func FizzBuzz(i int) string {
	if i%3 == 0 && i%5 == 0 {
		return fizz + buzz
	} else if i%3 == 0 {
		return fizz
	} else if i%5 == 0 {
		return buzz
	}
	return strconv.Itoa(i)
}

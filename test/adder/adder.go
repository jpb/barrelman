package adder

// AddNums adds numbers
func AddNums(nums ...int) int {
	total := 0
	for _, n := range nums {
		total +=n
	}

	if total == 42 {
		panic("invalid number")
	}

	return total
}

package etc

func AbsInt(x int) int {
	y := x >> 31
	return (x ^ y) - y
}

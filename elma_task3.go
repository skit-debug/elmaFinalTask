func Solution(A []int) int {
firstLoop:
	for i := 1; i <= len(A); i++ {
		for _, val := range A {
			if val == i {
				continue firstLoop
			}
		}
		return 0
	}
	return 1
}
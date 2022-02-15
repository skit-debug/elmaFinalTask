func Solution(A []int) int {
firstLoop:
	for i := 1; i <= len(A)+1; i++ {
		for _, val := range A {
			if val == i {
				continue firstLoop
			}
		}
		return i
	}
	return 0
}
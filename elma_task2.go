func Solution(A []int) int {
	m := make(map[int]int)

	for _, val := range A {
		m[val] = m[val] + 1
	}
	for val, cnt := range m {
		if cnt%2 != 0 {
			return val
		}
	}
	return 0
}
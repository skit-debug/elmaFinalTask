package solver

func Solution1(A []int, K int) []int {
	len := len(A)
	if len == 0 || K == 0 { //N=0 or K=0 case
		return A
	}
	a2 := make([]int, len)
	copy(a2, A)
	res := make([]int, len)
	for i := 0; i < K; i++ {
		res[0] = a2[len-1]
		copy(res[1:], a2[:len-1])
		copy(a2, res)
	}
	return res
}

func Solution2(A []int) int {
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

func Solution3(A []int) int {
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

func Solution4(A []int) int {
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

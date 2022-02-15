func Solution(A []int, K int) []int {
	len := len(A)
	if len == 0 || K == 0 { //N=0 or K=0 case
		return A
	}
	res := make([]int, len)
	for i := 0; i < K; i++ {
		res[0] = A[len-1]
		copy(res[1:], A[:len-1])
		copy(A, res)
	}
	return res
}
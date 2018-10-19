package algorithm

func diff2(nums []int, result *[2]int) {
	xor := 0
	for i := 0; i < len(nums); i++ {
		xor ^= nums[i]
	}
}

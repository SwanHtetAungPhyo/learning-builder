package main_test

import "testing"

func TwoSum(nums []int, target int) []int {
	var indices = make([]int, 2)
	for i := 0; i < len(nums); i++ {
		compliment := target - nums[i]
		for j := i + 1; j < len(nums); j++ {
			if nums[j] == compliment {
				indices[0] = i
				indices[1] = j
				return indices
			}
		}
	}
	return nil
}

func BenchmarkTwoSum(b *testing.B) {
	nums := []int{2, 7, 11, 15, 1, 8, 3, 6, 4, 5, 9, 10, 13, 14}
	target := 17
	for i := 0; i <= b.N; i++ {
		_ = TwoSum(nums, target)
	}
}

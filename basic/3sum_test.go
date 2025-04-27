package main_test

import (
	"sort"
	"testing"
)

func ThreeSum(nums []int) [][]int {
	var result [][]int
	nums = MergeSort(nums)

	for i := 0; i < len(nums); i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left, right := i+1, len(nums)-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}
				left++
				right--
			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	return result
}
func ThreeSum2(nums []int) [][]int {
	var result [][]int
	sort.Ints(nums)

	for i := 0; i < len(nums); i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left, right := i+1, len(nums)-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})
				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}
				left++
				right--
			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	return result
}
func MergeSort(nums []int) []int {
	if len(nums) <= 1 {
		return nums
	}
	mid := len(nums) / 2
	left := MergeSort(nums[:mid])
	right := MergeSort(nums[mid:])
	return Merge(left, right)
}

func Merge(left []int, right []int) []int {
	result := []int{} // don't preallocate with length
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func TestThreeSum(t *testing.T) {
	testCases := []struct {
		nums   []int
		expect [][]int
	}{
		{[]int{-1, 0, 1, 2, -1, -4}, [][]int{{-1, -1, 2}, {-1, 0, 1}}},
		{[]int{0, 0, 0, 0}, [][]int{{0, 0, 0}}},
		{[]int{1, 2, -2, -1}, [][]int{}},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			got := ThreeSum(tc.nums)
			if len(got) != len(tc.expect) {
				t.Errorf("expected %v, got %v", tc.expect, got)
			}
			for i := range got {
				for j := range got[i] {
					if got[i][j] != tc.expect[i][j] {
						t.Errorf("expected %v, got %v", tc.expect, got)
					}
				}
			}
		})
	}
}

func BenchmarkThreeSum(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ThreeSum([]int{-1, 0, 1, 2, -1, -4})
	}
}

func BenchmarkThreeSum2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ThreeSum2([]int{-1, 0, 1, 2, -1, -4})
	}
}

package main_test

import "testing"

func ContainDuplicate(nums []int) bool {
	var hashMap = make(map[int]struct{}, len(nums))
	if len(nums) == 0 {
		return false
	}
	for _, num := range nums {
		if _, exist := hashMap[num]; exist {
			return true
		}
		hashMap[num] = struct{}{}
	}
	return false
}

func BenchmarkContainDuplicate(b *testing.B) {
	nums := []int{1, 2, 3, 1}
	for i := 0; i <= b.N; i++ {
		_ = ContainDuplicate(nums)
	}
}
func TestContainDuplicate(t *testing.T) {
	tests := []struct {
		nums     []int
		expected bool
	}{
		{[]int{1, 2, 3, 1}, true},
		{[]int{1, 2, 3, 4}, false},
		{[]int{}, false},
		{[]int{1, 1, 1, 1}, true},
	}

	for _, test := range tests {
		actual := ContainDuplicate(test.nums)
		if actual != test.expected {
			t.Errorf("ContainDuplicate(%v) = %v, want %v", test.nums, actual, test.expected)
		}
	}

}

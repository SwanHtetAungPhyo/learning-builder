package main

import "fmt"

func TwoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		compliment := target - nums[i]
		for j := i + 1; j < len(nums); j++ {
			if nums[j] == compliment {
				return []int{i, j}
			}
		}
	}
	return nil
}

func main() {
	fmt.Println(TwoSum([]int{2, 7, 11, 15}, 9))

}

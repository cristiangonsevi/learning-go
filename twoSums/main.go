package main

import (
	"fmt"
)

func main() {
	arrNums := []int{3, 2, 4}
	fmt.Println(twoSum(arrNums, 6))
}

func twoSum(nums []int, target int) []int {
	dicc := make(map[int]int)

	for i, value := range nums {
		complement := target - value

		if idx, ok := dicc[complement]; ok {
			return []int{idx, i}
		}

		dicc[value] = i

	}

	return []int{}
}

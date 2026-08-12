package main

func sumEven(nums []int) int {
	var acc int
	for _, num := range nums {
		if num%2 == 0 {
			acc += num
		}
	}
	return acc
}

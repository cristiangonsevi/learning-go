package main

import "fmt"

func main() {
	result := sum(1, 2)
	fmt.Println("Result of sum is: ", result)
}

func sum(a int, b int) int {
	return a + b
}

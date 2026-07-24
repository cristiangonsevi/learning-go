package main

import "fmt"

func main() {
	fmt.Println("Hello, World1")
	var x = 5

	y := &x

	*y = 10

	fmt.Println("Values of %s %#v", &y, *y)
}

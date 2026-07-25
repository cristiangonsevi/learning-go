package main

import (
	"fmt"
	"listing/pkg/utils"
)

func main() {
	fmt.Println("Hello, World")
	err := utils.ListCurrentDirectory()

	if err != nil {
		fmt.Println("Error listando directorio")
	}
}

package main

import (
	"fmt"
	"week1/utils"
)

func main() {
	var n1, n2 float64
	fmt.Println("Ingrese el primer numero decimal...")
	fmt.Scanln(&n1)
	fmt.Println("Ingrese el segundo numero decimal...")
	fmt.Scanln(&n2)
	result := utils.DecimalToMilliInt(n1) + utils.DecimalToMilliInt(n2)
	fmt.Println(utils.MilliIntToDecimal(result))
}

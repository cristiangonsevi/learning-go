package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parseToNumber(s string) (int, error) {
	n, err := strconv.Atoi(s)

	if err != nil {
		return 0, errors.New("No es un numero entero valido")
	}
	return n, nil
}

func main() {
	var number string
	fmt.Println("Ingrese un numero entero....")
	fmt.Scanln(&number)
	result, err := parseToNumber(number)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("El numero ingresado es:", result)
}

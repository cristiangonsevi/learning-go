package main

import "fmt"

func main() {
	nums := []int{1, 2, 3, 4, 5, 6}
	result := sumEven(nums)

	fmt.Printf("Resultado de sumar solo numeros pares %v de array de numeros %v\n", result, nums)

	vowels := "hola"

	result = countVowels(vowels)

	fmt.Printf("Total de vocales %v en texto %v\n", result, vowels)
}

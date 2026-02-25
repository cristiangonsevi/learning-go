package main

import (
	"fmt"
	"reflect"
)

type Account struct {
	balance float64
}

func (a *Account) Deposit(amount float64) {
	varType := reflect.TypeOf(amount)
	if varType.Kind() != reflect.Float64 {
		fmt.Println("Error: El monto debe ser un número decimal.")
		return
	}
	fmt.Println("Tipo de variable:", varType)
	a.balance += amount
}

func (a *Account) Balance() float64 {
	return a.balance
}

func main() {
	acc := Account{}
	acc.Deposit(1.1)
	fmt.Println(acc.Balance())
	acc.Deposit(1.3)
	fmt.Println(acc.Balance())
}

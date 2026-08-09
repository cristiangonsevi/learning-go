package main

import (
	"fmt"
	"sync"
)

type Account struct {
	Name    string
	Balance float64
}

func bank() {
	var wg sync.WaitGroup
	var mu sync.Mutex

	myAccount := Account{Name: "My Account", Balance: 1000}

	for range 500 {
		wg.Go(func() {
			mu.Lock()
			defer mu.Unlock()
			addBalance(&myAccount, 10)
		})
	}

	for range 300 {
		wg.Go(func() {
			mu.Lock()
			defer mu.Unlock()
			withdrawBalance(&myAccount, 5)
		})
	}

	wg.Wait()

	fmt.Println("Account ", myAccount)
}

func addBalance(account *Account, amount float64) {
	account.Balance += amount
}

func withdrawBalance(account *Account, amount float64) {
	if account.Balance >= amount {
		account.Balance -= amount
	}
}

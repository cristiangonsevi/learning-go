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

	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			addBalance(&myAccount, 10)
		}()
	}

	for i := 0; i < 300; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			withdrawBalance(&myAccount, 5)
		}()
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

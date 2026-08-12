package main

import (
	"fmt"
	"sync"
)

type User struct {
	Name  string
	Email string
}

func bootstrap() {
	var wg sync.WaitGroup
	emailChan := make(chan User, 100)
	for i := range 5 {
		wg.Add(1)
		workerID := i + 1
		go func(id int) {
			emailConsumer(&wg, emailChan, id)
		}(workerID)
	}
	go func() {
		emailProducer(emailChan)
	}()

	wg.Wait()
}

func emailProducer(emailChan chan<- User) {
	for i := range 100 {
		name := fmt.Sprintf("Usuario %v", i+1)
		email := fmt.Sprintf("useremail%d@mail.com", i+1)
		emailChan <- User{
			Name:  name,
			Email: email,
		}
	}
	close(emailChan)
}

func emailConsumer(wg *sync.WaitGroup, emailChan <-chan User, workerID int) {
	defer wg.Done()
	for emailC := range emailChan {
		fmt.Printf("Procesando email %s en el workerId %d\n", emailC.Email, workerID)
	}
}

package main

import (
	"fmt"
	"sync"
)

type Order struct {
	ID int
}

func workerPool() {
	var wg sync.WaitGroup
	ordersChan := make(chan Order, 100)

	for i := range 5 {
		wg.Add(1)
		workerID := i + 1
		go func(ordersChan <-chan Order) {
			ordersConsumer(&wg, ordersChan, workerID)
		}(ordersChan)
	}

	ordersProducer(ordersChan)

	close(ordersChan)
	wg.Wait()
}

func ordersConsumer(wg *sync.WaitGroup, orderChan <-chan Order, worker int) {
	defer wg.Done()
	for order := range orderChan {
		fmt.Printf("Worker #%v orden #%v\n", worker, order.ID)
	}
}

func ordersProducer(ordersChan chan<- Order) {
	for i := range 100 {
		ordersChan <- Order{ID: i + 1}
	}
}

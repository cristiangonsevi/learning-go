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

	for i := 0; i < 5; i++ {
		wg.Add(1)
		workerId := i + 1
		go func(ordersChan <-chan Order) {
			ordersConsumer(&wg, ordersChan, workerId)
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
	for i := 0; i < 100; i++ {
		ordersChan <- Order{ID: i + 1}
	}

}

package main

import (
	"fmt"
	"sync"
)

func gorutines() {
	var wg sync.WaitGroup

	var mu sync.Mutex
	counter := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			defer mu.Unlock()
		}()
	}

	wg.Wait()

	fmt.Println("Counter ", counter)
}

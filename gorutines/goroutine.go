package main

import (
	"fmt"
	"sync"
)

func gorutines() {
	var wg sync.WaitGroup

	var mu sync.Mutex
	counter := 0

	for range 1000 {
		wg.Go(func() {
			mu.Lock()
			counter++
			defer mu.Unlock()
		})
	}

	wg.Wait()

	fmt.Println("Counter ", counter)
}

package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello world")
	begin := time.Now()
	fetchUrls()
	duration := time.Since(begin)

	fmt.Println("La tarea duro ", duration)
}

func fetchUrls() {
	var wg sync.WaitGroup
	urls := []string{"https://crisego.com", "https://termisearch.com"}

	wg.Add(len(urls))

	for _, url := range urls {
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)

			if err != nil {
				return
			}

			fmt.Printf("La pagina %v respondio con codigo %d\n", url, resp.StatusCode)
			resp.Body.Close()

		}(url)

	}

	wg.Wait()
}

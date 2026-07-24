package main

import (
	"fmt"
	"os"
	"sync"
)

func main() {
	fmt.Println("Listing current directory!")

	listDirectory()

}

func listDirectory() {
	var wg sync.WaitGroup

	files, err := os.ReadDir(".")

	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	wg.Add(len(files))

	for _, file := range files {
		go func(file os.DirEntry) {
			defer wg.Done()
			var emoji rune

			if file.IsDir() {
				emoji = '📁'
			} else {
				emoji = '🧾'
			}
			fmt.Printf("%c %v - con gorutines\n", emoji, file.Name())
			// time.Sleep(1 * time.Second)
		}(file)
	}

	wg.Wait()
}

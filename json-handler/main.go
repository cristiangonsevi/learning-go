package main

import (
	"encoding/json"
	"log"
	"os"
)

type Person struct {
	Name string `json:"name"`
}

func main() {
	file, err := os.OpenFile("./file.json", os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		log.Fatal("Error: ", err)
	}

	defer file.Close()

	var person Person

	json.NewDecoder(file).Decode(&person)
	output, _ := json.Marshal(person)

	log.Println("data", string(output))
}

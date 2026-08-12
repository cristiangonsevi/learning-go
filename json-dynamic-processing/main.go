package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	jsonData := `{"name": "Cristian", "age": 30}`

	data := make(map[string]interface{})

	json.Unmarshal([]byte(jsonData), &data)

	fmt.Printf("Valor de jsonData %v es %v", jsonData, data)
}

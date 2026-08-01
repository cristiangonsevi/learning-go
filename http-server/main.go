package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthStatus struct {
	Status string `json:"status"`
}

func health(w http.ResponseWriter, req *http.Request) {

	response := HealthStatus{Status: "OK"}
	output, err := json.Marshal(response)

	if err != nil {
		http.Error(w, "Error interno", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(output)

}

func main() {
	http.HandleFunc("/health", health)

	log.Println("Servidor corriendo en puerto 8080")
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal(err)
	}

}

package main

import (
	"log"
	"math/rand"
	"net/http"
	"url-shortener/internal/config"

	"github.com/gin-gonic/gin"
)

const BASE_URL = "crisego.com"

func main() {

	config, err := config.Load()

	if err != nil {
		log.Fatal("Error loading env configuration", err)
	}

	log.Println("configuration ", config)

	router := gin.Default()

	router.GET("/short", generateShortUrl)

	http.ListenAndServe(":8080", router)
}

func generateShortUrl(c *gin.Context) {
	url := BASE_URL + "/" + generateCode()

	c.JSON(200, gin.H{
		"url": url,
	})
}

func generateCode() string {
	randomString := make([]byte, 6)

	rangeString := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	for i := 0; i < 6; i++ {
		randomString[i] = rangeString[rand.Intn(len(rangeString))]
	}

	return string(randomString)
}

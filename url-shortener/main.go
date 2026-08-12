package main

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

const BASE_URL = "crisego.com"

func main() {

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

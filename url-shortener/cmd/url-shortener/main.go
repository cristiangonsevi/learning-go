package main

import (
	"fmt"
	"log"
	"math/rand"
	"url-shortener/internal/config"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository/postgres"
	"url-shortener/internal/router"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {

	config, err := config.Load()

	userRepo, _ := postgres.New(config.DBName)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewAuthHandler(userService)

	fmt.Println(userHandler)

	if err != nil {
		log.Fatal("Error loading env configuration", err)
	}

	r := router.New(userHandler)

	r.Run(":8080")

}

func generateShortUrl(c *gin.Context) {
	url := "/" + generateCode()

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

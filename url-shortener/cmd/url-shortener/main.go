package main

import (
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

	cfg, err := config.Load()

	if err != nil {
		log.Fatal("Error loading env vars: ", err)
	}

	connStr := config.ConnectionString(cfg)

	repo, err := postgres.New(connStr)

	if err != nil {
		log.Fatal("Error initializating repository")
	}

	err = repo.TestConnection()

	if err != nil {
		log.Fatal("Error connecting to database", err)
	}

	log.Println("[DATABASE] Online")

	userService := services.NewUserService(repo)
	userHandler := handler.NewAuthHandler(userService)

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

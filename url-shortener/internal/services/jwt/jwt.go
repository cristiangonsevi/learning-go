package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(data map[string]interface{}) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"exp": time.Now().Add(15 * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}

	for k, v := range data {
		claims[k] = v
	}

	fmt.Println("data ", claims, secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return tokenStr, nil

}

func GenerateRefreshToken(payload map[string]interface{}, expiresAt time.Time) (string, error) {
	payload["iat"] = expiresAt
	token, err := GenerateJWT(payload)
	if err != nil {
		return "", err
	}
	return token, nil
}

func validateJWT() {}

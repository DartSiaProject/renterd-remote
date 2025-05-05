package auth

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtSecretKey is the secret key used to sign the JWT tokens. It can be set via the JWT_SECRET environment variable.
// If not set, it defaults to "default
var JwtSecretKey = []byte(func() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "1hwGE8Y6nHbPVRA9"
}())

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": os.Getenv("USER_EMAIL"),
			//Jeton validty : 7 days
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

	tokenString, err := token.SignedString(JwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return JwtSecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

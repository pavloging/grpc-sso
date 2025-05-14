package jwt

import (
	"grpc-sso/internal/domain/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(duration).Unix()

	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	secretJwtKey := os.Getenv("SECRET_JWT_KEY")
	if secretJwtKey == "" {
		panic("SECRET_JWT_KEY is not set")
	}

	tokenString, err := token.SignedString([]byte(secretJwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

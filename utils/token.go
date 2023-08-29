package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type TokenPayload struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func GetToken(authHeader string) string {
	var token string
	fields := strings.Fields(authHeader)

	if len(fields) > 1 && fields[0] == "Bearer" {
		token = fields[1]
	}

	if token == "" {
		return "You need to login first"
	}

	return token
}

func ValidateToken(token string, secret string) (TokenPayload, error) {

	claims := jwt.MapClaims{}

	tok, err := jwt.ParseWithClaims(token, claims, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil || !tok.Valid {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}

	userIDStr, ok := claims["ID"].(string)
	if !ok {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}

	ID, err := uuid.Parse(userIDStr)
	if err != nil {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}

	return TokenPayload{
		ID:    ID,
		Email: email,
	}, nil
}

func VerifyToken(c *fiber.Ctx) (string, *fiber.Error) {

	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "No authorization header found")
	}
	token := GetToken(authHeader)

	if token == "" {
		return "", fiber.NewError(fiber.StatusBadRequest, "No token found")
	}

	res, err := ValidateToken(token, viper.GetString("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return res.Email, nil
}

func GenerateToken(userID uuid.UUID, email string, secret string, expiry time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":    userID.String(),
		"email": email,
		"exp":   time.Now().Add(expiry).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uuid.UUID, email string) (string, error) {
	return GenerateToken(userID, email, viper.GetString("REFRESH_TOKEN_SECRET"), viper.GetDuration("REFRESH_TOKEN_EXPIRY"))
}

func GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return GenerateToken(userID, email, viper.GetString("ACCESS_TOKEN_SECRET"), viper.GetDuration("ACCESS_TOKEN_EXPIRY"))
}

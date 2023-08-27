package utils

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

	role, ok := claims["role"].(string)
	if !ok {
		return TokenPayload{}, fmt.Errorf("invalid token")
	}

	return TokenPayload{
		ID:    ID,
		Email: email,
		Role:  role,
	}, nil
}

package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length")
	}

	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	max := big.NewInt(int64(len(charset)))

	otp := make([]byte, length)
	for i := range otp {
		index, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		otp[i] = charset[index.Int64()]
	}

	return string(otp), nil
}

package utils

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var LENGTH = big.NewInt(int64(len(ALPHABET)))

func GenerateToken(length int) string {
	var builder strings.Builder
	builder.Grow(length)
	for i := 0; i < length; i++ {
		value, err := rand.Int(rand.Reader, LENGTH)
		if err != nil {
			panic(err)
		}
		builder.WriteByte(ALPHABET[value.Int64()])
	}
	return builder.String()
}

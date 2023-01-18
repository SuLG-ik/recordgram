package utils

import (
	"crypto/rand"
	"math/big"
	"recordgram/config"
	"regexp"
	"strconv"
	"strings"
)

const ALPHABET = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var LENGTH = big.NewInt(int64(len(ALPHABET)))

func GenerateTokenWithLength(length int) string {
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

func GenerateTokenFromConfig(config config.Config) string {
	return GenerateTokenWithLength(config.Api.KeyLength)
}

var tokenRegex *regexp.Regexp

func IsTokenValid(token string, config config.Config) bool {
	if tokenRegex == nil {
		tokenRegex, _ = regexp.Compile("\\d{1,9}:[a-zA-Z0-9]{" + strconv.Itoa(config.Api.KeyLength) + "}")
	}
	return tokenRegex.MatchString(token)
}

package passwords

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func HashPassword(rawPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	return "{bcrypt}" + string(bytes), err
}

func MatchPassword(rawPassword, hash string) bool {
	if strings.HasPrefix(hash, "{bcrypt}") {
		password := strings.TrimPrefix(hash, "{bcrypt}")
		err := bcrypt.CompareHashAndPassword([]byte(password), []byte(rawPassword))
		return err == nil
	}
	log.Panicf("Unknown hashing method")
	return false
}

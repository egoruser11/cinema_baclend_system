package utils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
)

const DbConfig = "host=localhost user=user password=password dbname=mydb port=5432 sslmode=disable"

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxy" +
		"0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func InputErrorsValid(errorsValid map[string]string) string {
	errorsRes := []string{}
	for field, err := range errorsValid {
		errorsRes = append(errorsRes, field+" :"+err)
	}
	return strings.Join(errorsRes, "\n")
}

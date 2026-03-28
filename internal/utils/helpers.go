package utils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
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

func ParseSeats(seatsStr string) map[uint][]uint {
	result := make(map[uint][]uint)
	seatsStr = strings.TrimSuffix(seatsStr, ",")
	if seatsStr == "" {
		return result
	}
	seatsArr := strings.Split(seatsStr, ",")

	for _, seat := range seatsArr {
		parts := strings.Split(strings.TrimSpace(seat), " - ")
		if len(parts) != 2 {
			continue
		}

		row, err := strconv.ParseUint(parts[0], 10, 0)
		if err != nil {
			continue
		}

		num, err := strconv.ParseUint(parts[1], 10, 0)
		if err != nil {
			continue
		}

		rowUint := uint(row)
		numUint := uint(num)

		exists := false
		for _, existing := range result[rowUint] {
			if existing == numUint {
				exists = true
				break
			}
		}

		if !exists {
			result[rowUint] = append(result[rowUint], numUint)
		}
	}

	return result
}

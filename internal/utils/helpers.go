package utils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

var badWords = []string{
	"fuck", "shit", "damn", "bitch", "cunt", "dick", "cock",
	"pussy", "asshole", "bastard", "whore", "slut", "cocksucker",
	"motherfucker", "fck", "fk", "sh1t", "b1tch", "mf", "mfucker",
	"retard", "moron", "rape", "molest", "incest", "pedo", "pervert",
}

var punctuationRegex = regexp.MustCompile(`[^\w\s]`)

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

func ValidateCommentInReview(comment string) bool {
	cleaned := punctuationRegex.ReplaceAllString(comment, "")
	words := strings.Fields(cleaned)

	for _, word := range words {
		wordLower := strings.ToLower(word)
		wordRunes := []rune(wordLower)

		if len(wordRunes) > 50 {
			return false
		}
		for _, badWord := range badWords {
			if removeRepeatedChars(wordRunes) == badWord {
				return false
			}
		}

		if isLink(wordLower) {
			return false
		}

		if isEmail(wordLower) {
			return false
		}

		if isPhoneNumber(wordLower) {
			return false
		}

		if isIPAddress(wordLower) {
			return false
		}
	}
	return true
}

func isLink(word string) bool {
	domains := []string{".com", ".ru", ".net", ".org", ".io", ".xyz", ".site", ".online", ".club"}
	for _, domain := range domains {
		if strings.Contains(word, domain) {
			return true
		}
	}

	if strings.HasPrefix(word, "http") || strings.HasPrefix(word, "www") {
		return true
	}

	linkPattern := regexp.MustCompile(`^[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}$`)
	return linkPattern.MatchString(word)
}

func isEmail(word string) bool {
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailPattern.MatchString(word)
}

func isPhoneNumber(word string) bool {
	digits := regexp.MustCompile(`\d`).FindAllString(word, -1)
	if len(digits) >= 7 && len(digits) <= 15 {
		return true
	}
	return false
}

func isIPAddress(word string) bool {
	ipPattern := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	return ipPattern.MatchString(word)
}

func removeRepeatedChars(runes []rune) string {
	if len(runes) == 0 {
		return ""
	}
	result := []rune{runes[0]}
	for i := 1; i < len(runes); i++ {
		if runes[i] != result[len(result)-1] {
			result = append(result, runes[i])
		}
	}
	return string(result)
}

func GetSLiceSum(slice []float64) float64 {
	result := 0
	for _, v := range slice {
		result += int(v)
	}
	return float64(result)
}

package validators

import (
	"gorm.io/gorm"
	"regexp"
	"strings"
)

type RegisterRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Age        uint   `json:"age"`
	DeviceInfo string `json:"device_info"`
}

func ValidateRegister(db *gorm.DB, req RegisterRequest) (map[string]string, bool) {

	errors := make(map[string]string)
	if req.Username == "" {
		errors["Username"] = "Username is required"
	} else if len(req.Username) < 3 {
		errors["username"] = "Username must be at least 3 characters"
	} else if len(req.Username) > 50 {
		errors["username"] = "Username must be less than 50 characters"
	} else if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(req.Username) {
		errors["username"] = "Username can only contain letters, numbers, dots and underscores"
	}

	if req.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(req.Email) {
		errors["email"] = "Invalid email"
	}

	if len(errors) == 0 {
		var exists bool
		db.Raw(`
			SELECT EXISTS(
				SELECT 1 FROM users 
				WHERE username = ? OR email = ?
			)
		`, req.Username, req.Email).Scan(&exists)

		if exists {
			errors["username"] = "Username or email already exists"
		}
	}

	return errors, len(errors) == 0

}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(strings.ToLower(email))
}

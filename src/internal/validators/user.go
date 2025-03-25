package validators

import (
	"unicode"
)

func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)
	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}
	return hasUpperCase && hasLowerCase && hasDigit && hasSpecial
}

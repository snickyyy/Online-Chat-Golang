package validators

import (
	"github.com/go-playground/validator/v10"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"unicode"
)

func ValidatePassword(field validator.FieldLevel) bool {
	password := field.Field().String()
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

func ValidatorImage(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	ext := filepath.Ext(file.Filename)
	allowedExtensions := map[string]bool{".jpg": true, ".jpeg": true, ".png": true}
	if !allowedExtensions[ext] {
		return false
	}

	const maxSize = 8 << 20
	if file.Size > maxSize {
		return false
	}

	return true
}

func ValidateUsername(field validator.FieldLevel) bool {
	if len(field.Field().String()) < 1 || len(field.Field().String()) > 35 {
		return false
	}
	regex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return regex.MatchString(field.Field().String())
}

package validate

import (
	"net/mail"
	"unicode"

	"github.com/Artexus/api-widyabhuvana/src/constant"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return constant.ErrInvalid
	}
	return nil
}

func ValidatePassword(password string) error {
	isDigit := false
	isUpper := false
	for _, s := range password {
		if unicode.IsDigit(s) {
			isDigit = true
		} else if unicode.IsUpper(s) {
			isUpper = true
		}
	}

	if len(password) > constant.PasswordMaxlength ||
		!isDigit || !isUpper {
		return constant.ErrInvalid
	}
	return nil
}

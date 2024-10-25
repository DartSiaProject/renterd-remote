package utils

import (
	"errors"
	"net/mail"
	constants "renterd-remote/constant"
	"unicode"
)

func ValidateEmail(email interface{}) error {
	if _, ok := email.(string); !ok {
		return errors.New(constants.EmailError)
	}

	str, _ := email.(string)
	_, err := mail.ParseAddress(str)
	if err != nil {
		return errors.New(constants.EmailError)
	}
	return nil
}

func ValidatePassword(password interface{}) error {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	if _, ok := password.(string); !ok {
		return errors.New(constants.PasswordError)
	}

	str, _ := password.(string)

	if len(str) >= 8 {
		hasMinLen = true
	}
	for _, char := range str {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial {
		return nil
	} else {
		return errors.New(constants.PasswordError)
	}

}

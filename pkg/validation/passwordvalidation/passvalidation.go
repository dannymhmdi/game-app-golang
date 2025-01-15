package passwordvalidation

import (
	"errors"
	"unicode"
)

func IsPasswordValid(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var hasDigit, hasUpperCase, hasSpecialChar bool
	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpperCase = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecialChar = true
		}
	}

	if !hasDigit || !hasUpperCase || !hasSpecialChar {
		var err string
		if !hasDigit {
			err += "password must contain at least one digit\n"
		}

		if !hasUpperCase {
			err += "password must contain at least one upper case\n"
		}

		if !hasSpecialChar {
			err += "password must contain at least one special character\n"
		}

		return errors.New(err)
	}

	return nil
}

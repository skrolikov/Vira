package validators

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidateCredentials проверяет имя пользователя и пароль на соответствие требованиям
func ValidateCredentials(username, password string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}

	return ValidatePassword(password)
}

// ValidateUsername проверяет имя пользователя
func ValidateUsername(username string) error {
	if strings.TrimSpace(username) == "" {
		return errors.New("имя пользователя не может быть пустым")
	}

	if !usernameRegex.MatchString(username) {
		return errors.New("имя пользователя должно содержать только латинские буквы, цифры и подчёркивания (3–32 символа)")
	}

	return nil
}

// ValidatePassword проверяет пароль на сложность
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("пароль должен содержать не менее 8 символов")
	}

	if len(password) > 64 {
		return errors.New("пароль не должен превышать 64 символа")
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasDigit   = false
		hasSpecial = false
	)

	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	// Требования к паролю:
	// - минимум 1 заглавная буква
	// - минимум 1 строчная буква
	// - минимум 1 цифра
	// - минимум 1 спецсимвол
	if !hasUpper {
		return errors.New("пароль должен содержать хотя бы одну заглавную букву")
	}

	if !hasLower {
		return errors.New("пароль должен содержать хотя бы одну строчную букву")
	}

	if !hasDigit {
		return errors.New("пароль должен содержать хотя бы одну цифру")
	}

	if !hasSpecial {
		return errors.New("пароль должен содержать хотя бы один специальный символ (!@#$%^&* и т.д.)")
	}

	return nil
}

// ValidateEmail проверяет email на корректность
func ValidateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email не может быть пустым")
	}

	if len(email) > 254 {
		return errors.New("email не должен превышать 254 символа")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("неверный формат email")
	}

	// Дополнительная проверка доменной части
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("неверный формат email")
	}

	domain := parts[1]
	if strings.Contains(domain, "..") {
		return errors.New("неверный формат домена в email")
	}

	return nil
}

// ValidateUserInput проверяет все данные пользователя
func ValidateUserInput(username, password, email string) error {
	if err := ValidateUsername(username); err != nil {
		return err
	}

	if err := ValidatePassword(password); err != nil {
		return err
	}

	if email != "" {
		if err := ValidateEmail(email); err != nil {
			return err
		}
	}

	return nil
}

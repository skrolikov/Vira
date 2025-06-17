package service

import (
	"errors"
	"regexp"
	"strings"
	"vira-id/internal/types"
)

func (s *AuthService) validateRegistration(req types.RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("имя пользователя обязательно")
	}
	if len(req.Password) < 6 {
		return errors.New("пароль должен быть не менее 6 символов")
	}
	if req.Email != "" {
		// простой regex для проверки email (можно сделать лучше)
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			return errors.New("некорректный email")
		}
	}
	return nil
}

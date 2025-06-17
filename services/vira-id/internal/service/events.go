package service

import (
	"context"
	"time"
	"vira-id/internal/events"
)

// emitRegistrationEvent отправляет событие регистрации пользователя в Kafka.
// Используется context с таймаутом 3 секунды, чтобы не блокировать навсегда
// в случае проблем с отправкой.
// userID, username — данные пользователя,
// ip — IP адрес,
// device — User-Agent или описание устройства,
// время события берётся в момент вызова.
func (s *AuthService) emitRegistrationEvent(ctx context.Context, userID, username, ip, device string) {
	// Создаём контекст с таймаутом 3 секунды
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Вызываем функцию отправки события регистрации
	err := events.EmitUserRegisteredEvent(ctxTimeout, s.Producer, s.Logger, userID, username, ip, device, time.Now())
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение (fire-and-forget)
		s.Logger.Error("Ошибка отправки события регистрации: %v", err)
	}
}

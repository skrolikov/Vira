package handlers

import (
	"encoding/json"
	"net/http"

	"vira-id/internal/service"
	"vira-id/internal/types"
)

// RegisterHandler обрабатывает регистрацию нового пользователя
// Добавлены kafka.Producer и logger для отправки Kafka-события
func RegisterHandler(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		userAgent := r.UserAgent()
		ip := getIP(r)

		resp, err := authService.Register(r.Context(), req, ip, userAgent)
		if err != nil {
			// Обработка ошибок и ответ с нужным статусом
			http.Error(w, err.Error(), http.StatusBadRequest) // пример
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

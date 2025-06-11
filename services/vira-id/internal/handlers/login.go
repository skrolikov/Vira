package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"vira-id/internal/service"
	"vira-id/internal/types"
)

// LoginHandler обрабатывает авторизацию пользователя
func LoginHandler(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
			return
		}

		req.Username = strings.TrimSpace(req.Username)
		req.Password = strings.TrimSpace(req.Password)

		if len(req.Username) < 3 || len(req.Password) < 6 {
			http.Error(w, "Неверный логин или пароль", http.StatusBadRequest)
			return
		}

		userAgent := r.UserAgent()
		ip := getIP(r)

		resp, err := authService.Login(r.Context(), req, ip, userAgent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func getIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

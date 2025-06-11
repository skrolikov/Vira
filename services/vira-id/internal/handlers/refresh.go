package handlers

import (
	"encoding/json"
	"net/http"

	"vira-id/internal/service"
	"vira-id/internal/types"
)

func RefreshHandler(svc *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefreshRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		tokens, err := svc.RefreshToken(r.Context(), req.RefreshToken, r.RemoteAddr, r.UserAgent())
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"access_token":  tokens.Access,
			"refresh_token": tokens.Refresh,
		})
	}
}

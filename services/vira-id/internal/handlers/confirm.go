package handlers

import (
	"net/http"

	"vira-id/internal/service"
)

func ConfirmUserHandler(authService *service.AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		email := r.URL.Query().Get("email")
		token := r.URL.Query().Get("token")

		if email == "" || token == "" {
			http.Error(w, "необходимо указать email и токен", http.StatusBadRequest)
			return
		}

		err := authService.ConfirmUser(ctx, email, token)
		if err != nil {
			http.Error(w, "не удалось подтвердить пользователя: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Пользователь успешно подтверждён"))
	}
}

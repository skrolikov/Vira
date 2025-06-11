package handlers

import (
	"encoding/json"
	"net/http"
	"vira-id/internal/types"

	middleware "github.com/skrolikov/vira-middleware"

	db "github.com/skrolikov/vira-db"
)

func MeHandler(repo db.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		if userID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := repo.GetUserByID(userID)
		if err != nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(types.UserInfo{
			ID:       user.ID,
			Username: user.Username,
		})
	}
}

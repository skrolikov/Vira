package handlers

import (
	"net/http"

	middleware "github.com/skrolikov/vira-middleware"

	"github.com/gorilla/mux"

	"log"

	"github.com/redis/go-redis/v9"
	config "github.com/skrolikov/vira-config"
)

func DeleteSessionHandler(cfg *config.Config, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		vars := mux.Vars(r)
		sessionID := vars["id"]
		if sessionID == "" {
			http.Error(w, "Неверный ID сессии", http.StatusBadRequest)
			return
		}

		// Удаляем сессию из Redis
		sessionKey := "session:" + userID + ":" + sessionID
		if err := rdb.Del(r.Context(), sessionKey).Err(); err != nil {
			log.Printf("Ошибка удаления сессии: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// Если этот refresh токен совпадает с основным в ключе refresh:<userID>, удаляем и его (если надо)
		key := "refresh:" + userID
		savedToken, err := rdb.Get(r.Context(), key).Result()
		if err == nil && savedToken == sessionID {
			if err := rdb.Del(r.Context(), key).Err(); err != nil {
				log.Printf("Ошибка удаления refresh токена: %v", err)
			}
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

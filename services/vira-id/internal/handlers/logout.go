package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"vira-id/internal/types"

	"log"

	"github.com/redis/go-redis/v9"
	config "github.com/skrolikov/vira-config"
	jwt "github.com/skrolikov/vira-jwt"
)

func LogoutHandler(cfg *config.Config, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LogoutRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.RefreshToken) == "" {
			http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
			return
		}

		claims, err := jwt.ParseToken(req.RefreshToken, cfg.JwtSecret)
		if err != nil || !jwt.IsTokenType(claims, "refresh") {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok || userID == "" {
			http.Error(w, "Невалидный user_id", http.StatusUnauthorized)
			return
		}

		// Проверяем, совпадает ли токен с сохранённым refresh в Redis
		key := "refresh:" + userID
		savedToken, err := rdb.Get(r.Context(), key).Result()
		if err == redis.Nil {
			http.Error(w, "Токен не найден", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("Ошибка Redis при получении токена: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		if savedToken == req.RefreshToken {
			// Удаляем основной refresh токен (завершение сессии)
			if err := rdb.Del(r.Context(), key).Err(); err != nil {
				log.Printf("Ошибка удаления refresh токена: %v", err)
			}
		}

		// Удаляем сессию по ключу session:<userID>:<refresh>
		sessionKey := "session:" + userID + ":" + req.RefreshToken
		if err := rdb.Del(r.Context(), sessionKey).Err(); err != nil {
			log.Printf("Ошибка удаления сессии: %v", err)
		}

		w.WriteHeader(http.StatusNoContent) // 204 No Content
	}
}

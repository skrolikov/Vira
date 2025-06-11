package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"vira-id/internal/types"

	middleware "github.com/skrolikov/vira-middleware"

	"github.com/redis/go-redis/v9"
	config "github.com/skrolikov/vira-config"
)

func SessionsHandler(cfg *config.Config, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Авторизация: вытаскиваем userID из контекста, должен быть установлен middleware.Auth
		userID := middleware.GetUserID(r)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// 2. Парсим курсор из query-параметра. По умолчанию 0 (начало).
		cursorStr := r.URL.Query().Get("cursor")
		var cursor uint64
		if cursorStr != "" {
			var err error
			cursor, err = strconv.ParseUint(cursorStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid cursor", http.StatusBadRequest)
				return
			}
		}

		// 3. Готовим шаблон ключей: session:<userID>:*
		pattern := "session:" + userID + ":*"
		const countPerScan int64 = 20 // сколько ключей считываем за итерацию

		// 4. Запускаем SCAN: получаем срез ключей и новый курсор
		keys, nextCursor, err := rdb.Scan(r.Context(), cursor, pattern, countPerScan).Result()
		if err != nil {
			log.Printf("Ошибка SCAN Redis: %v", err)
			http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
			return
		}

		// 5. Для каждого ключа вытаскиваем JSON, парсим в types.SessionInfo
		sessions := make([]types.SessionInfo, 0, len(keys))
		for _, key := range keys {
			raw, err := rdb.Get(r.Context(), key).Result()
			if err != nil {
				log.Printf("Ошибка чтения сессии %s: %v", key, err)
				continue
			}

			var sess types.SessionInfo
			if err := json.Unmarshal([]byte(raw), &sess); err != nil {
				log.Printf("Ошибка парсинга сессии %s: %v", key, err)
				continue
			}

			// Извлекаем ID сессии из ключа "session:<userID>:<sessionID>"
			sess.ID = key[len("session:"+userID+":"):]
			sessions = append(sessions, sess)
		}

		// 6. Формируем ответ и отправляем
		resp := types.SessionsResponse{
			Cursor:   nextCursor,
			Sessions: sessions,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

package courses

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type Course struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var ctx = context.Background()

func GetCoursesHandler(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const cacheKey = "courses:all"

		// 1. Пытаемся получить из Redis
		cached, err := rdb.Get(ctx, cacheKey).Result()
		if err == nil {
			w.Header().Set("X-Cache", "HIT")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cached))
			return
		}

		// 2. Симулируем "из базы"
		courses := []Course{
			{ID: 1, Title: "Go для новичков"},
			{ID: 2, Title: "Docker и Kubernetes"},
		}

		body, _ := json.Marshal(courses)
		_ = rdb.Set(ctx, cacheKey, body, 5*time.Minute).Err() // кэш на 5 минут

		w.Header().Set("X-Cache", "MISS")
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}
}

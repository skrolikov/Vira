package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	middleware "github.com/skrolikov/vira-middleware"
)

func Proxy(target string) http.HandlerFunc {
	url, err := url.Parse(target)
	if err != nil {
		log.Printf("[Proxy] Ошибка парсинга target URL '%s': %v", target, err)
		return func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Неверный адрес прокси", http.StatusInternalServerError)
		}
	}

	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ModifyResponse = func(resp *http.Response) error {
		req := resp.Request
		if req != nil {
			log.Printf("[Proxy] Ответ от %s %s: статус %d", req.Method, req.URL.Path, resp.StatusCode)
		} else {
			log.Printf("[Proxy] Ответ с неизвестным запросом: статус %d", resp.StatusCode)
		}
		return nil
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Proxy] Входящий запрос: %s %s", r.Method, r.URL.String())

		userID := middleware.GetUserID(r)
		if userID != "" {
			r.Header.Set("X-User-ID", userID)
			log.Printf("[Proxy] Добавлен заголовок X-User-ID: %s", userID)
		}

		log.Printf("[Proxy] Проксирование на: %s%s", url.Host, r.URL.Path)

		proxy.ServeHTTP(w, r)
	}
}

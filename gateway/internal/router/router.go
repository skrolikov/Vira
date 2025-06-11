package router

import (
	"net/http"
	"vira-gateway/internal/proxy"

	"github.com/go-chi/chi/v5"
	config "github.com/skrolikov/vira-config"
	logger "github.com/skrolikov/vira-logger"
)

func Setup(cfg *config.Config, logger *logger.Logger) http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		r.Route("/id", func(r chi.Router) {
			r.Handle("/*", http.StripPrefix("/api/id", proxy.Proxy("http://vira-id:8080")))
		})

		r.Route("/dev", func(r chi.Router) {
			r.Handle("/*", http.StripPrefix("/api/dev", proxy.Proxy("http://vira-api-dev:8080")))
		})

		r.Route("/wish", func(r chi.Router) {
			r.Handle("/*", http.StripPrefix("/api/wish", proxy.Proxy("http://vira-api-wish:8080")))
		})
	})

	return r
}

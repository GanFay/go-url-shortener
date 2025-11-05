package httpx

import (
	"net/http"

	"github.com/GanFay/go-url-shortener/internal/config"
)

func NewRouter(cfg config.Config) *http.ServeMux {
	h := NewHandlers(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/version", h.Version)

	return mux
}

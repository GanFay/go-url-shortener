package httpx

import (
	"database/sql"
	"net/http"

	"github.com/GanFay/go-url-shortener/internal/config"
)

func NewRouter(cfg config.Config, db *sql.DB) *http.ServeMux {
	h := NewHandlers(cfg, db)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/version", h.Version)
	mux.HandleFunc("/db/ping", h.DBPing)
	mux.HandleFunc("/shorten", h.Shorten)
	mux.HandleFunc("/u/", h.Resolve)
	return mux
}

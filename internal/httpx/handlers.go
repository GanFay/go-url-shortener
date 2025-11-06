package httpx

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/GanFay/go-url-shortener/internal/config"
)

type Handlers struct {
	cfg config.Config
	db  *sql.DB
}

func NewHandlers(cfg config.Config, db *sql.DB) *Handlers { return &Handlers{cfg: cfg, db: db} }

func (h *Handlers) DBPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := h.db.Ping(); err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("pong"))
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, "ok")
	fmt.Println("ussed command /health")
}

func (h *Handlers) Version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintln(w, h.cfg.Version)
	fmt.Println("ussed command /version")
}

package httpx

import (
	"fmt"
	"net/http"

	"github.com/GanFay/go-url-shortener/internal/config"
)

type Handlers struct {
	cfg config.Config
}

func NewHandlers(cfg config.Config) *Handlers { return &Handlers{cfg: cfg} }

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

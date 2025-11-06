package httpx

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/GanFay/go-url-shortener/internal/config"
)

type Handlers struct {
	cfg config.Config
	db  *sql.DB
}

type shortenRequest struct {
	URL string `json:"url"`
}

func NewHandlers(cfg config.Config, db *sql.DB) *Handlers { return &Handlers{cfg: cfg, db: db} }

func (h *Handlers) Resolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	const prefix = "/u/"
	if !strings.HasPrefix(r.URL.Path, prefix) || len(r.URL.Path) <= len(prefix) {
		http.Error(w, "bad code", http.StatusBadRequest)
		return
	}
	code := r.URL.Path[len(prefix):]

	var url string
	err := h.db.QueryRow(`SELECT original_url FROM public.urls WHERE code = $1`, code).Scan(&url)
	if err == sql.ErrNoRows {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handlers) Shorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	var req shortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	code := generateCode(6)
	_, err = h.db.Exec(`INSERT INTO urls (code, original_url) VALUES ($1, $2)`, code, req.URL)
	if err != nil {
		http.Error(w, "failed to save", http.StatusInternalServerError)
		return
	}
	resp := map[string]string{"short_url": fmt.Sprintf("http://localhost:8080/u/%s", code)}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("failed to write response: %v\n", err)
		return
	}
}

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

func generateCode(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

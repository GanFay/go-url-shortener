package httpx

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

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

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	res, err := h.db.ExecContext(ctx, `UPDATE public.urls SET clicks = clicks + 1 WHERE code = $1`, code)
	if err != nil {
		slog.Error("failed to update clicks", "code", code, "err", err)
	} else {
		affected, _ := res.RowsAffected()
		if affected == 0 {
			slog.Warn("no rows updated (invalid code?)", "code", code)
		}
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.URL) == "" {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	var existingCode string
	err := h.db.QueryRow(`SELECT code FROM urls WHERE original_url = $1`, req.URL).Scan(&existingCode)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"url":"http://localhost:8080/u/%s","code":"%s","existing":true}`, existingCode, existingCode)
		return
	} else if err != sql.ErrNoRows {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for attmpt := 0; attmpt < 5; attmpt++ {
		code := generateCode(6)

		var insrt string
		err = h.db.QueryRow(
			`INSERT INTO urls (original_url, code)
					VALUES ($1, $2)
					ON CONFLICT (code) DO NOTHING
					RETURNING code`, req.URL, code).Scan(&insrt)

		if err == nil {
			short := fmt.Sprintf("http://localhost:8080/u/%s", insrt)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Location", short)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"url":"%s","code":"%s", "existing":false}`, short, insrt)
			return
		}
		if err.Error() != "sql: no rows in result set" {
			http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Error(w, "could not generate unique code", http.StatusConflict)
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

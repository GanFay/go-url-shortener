package httpx

import (
	"context"
	"encoding/hex"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func WithMiddleware(h http.Handler, mws ...func(http.Handler) http.Handler) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

type ctxKey string

const reqIDKey ctxKey = "req_id"

func RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := r.Header.Get("X-Request-ID")
			if id == "" {
				id = newReqID()
			}
			ctx := context.WithValue(r.Context(), reqIDKey, id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(sw, r)
			path := r.URL.Path
			if q := r.URL.RawQuery; q != "" {
				path += "?" + q
			}
			id, _ := r.Context().Value(reqIDKey).(string)
			slog.Info("http_request",
				"req_id", id,
				"method", r.Method,
				"path", path,
				"status", sw.status,
				"dur_ms", time.Since(start).Milliseconds(),
				"ua", r.UserAgent(),
				"ip", clientIP(r.Header.Get("X-Forwarded-For"), r.RemoteAddr),
			)
		})
	}
}
func clientIP(xff, remote string) string {
	if xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	return remote
}

func Recoverer() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					id, _ := r.Context().Value(reqIDKey).(string)
					slog.Error("panic", "req_id", id, "err", rec)
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func newReqID() string {
	var b [8]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

package main

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/GanFay/go-url-shortener/internal/config"
	"github.com/GanFay/go-url-shortener/internal/httpx"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}

	cfg := config.Load()
	db, err := sql.Open("pgx", cfg.DB_DSN)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		panic("db not reachable: " + err.Error())
	}
	mux := httpx.NewRouter(cfg, db)

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	handler := httpx.WithMiddleware(
		mux,
		httpx.Recoverer(),
		httpx.RequestID(),
		httpx.Logger(),
	)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

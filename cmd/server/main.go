package main

import (
	"database/sql"
	"fmt"
	"log"
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
	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr, "version:", cfg.Version)
	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("server error:", err)
	}
}

package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/GanFay/go-url-shortener/internal/config"
	"github.com/GanFay/go-url-shortener/internal/httpx"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
	}
	cfg := config.Load()

	mux := httpx.NewRouter(cfg)

	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr, "version:", cfg.Version)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("server error:", err)
	}
}

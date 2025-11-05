package main

import (
	"fmt"
	"net/http"

	"github.com/GanFay/go-url-shortener/internal/config"
	"github.com/GanFay/go-url-shortener/internal/httpx"
)

func main() {
	cfg := config.Load()

	mux := httpx.NewRouter(cfg)

	addr := ":" + cfg.Port
	fmt.Println("Listening on", addr, "version:", cfg.Version)

	if err := http.ListenAndServe(addr, mux); err != nil {
		fmt.Println("server error:", err)
	}
}

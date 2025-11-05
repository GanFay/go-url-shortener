package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/version", version)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("server error:", err)
		return
	}
	fmt.Println("server started on port 8080")
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = fmt.Fprintln(w, "ok")
	fmt.Println("used command /health")
}

func version(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	v := os.Getenv("APP_VERSION")
	if v == "" {
		v = "dev"
	}
	fmt.Fprintln(w, v)
}

package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", health)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "OK")
	if err != nil {
		fmt.Println(err)
		return
	}
}

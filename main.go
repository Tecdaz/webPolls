package main

import (
	"fmt"
	"net/http"
)

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	http.HandleFunc("/", serveIndex)

	port := ":8080"

	fmt.Println("Server started on port", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

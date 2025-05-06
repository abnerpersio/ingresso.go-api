package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"ingresso.go/config"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "pongs")
}

func main() {
	godotenv.Load()

	http.HandleFunc("/ping", pingHandler)

	port := config.GetEnv("PORT", "8080")
	addr := ":" + port

	fmt.Printf("Server starting on http://localhost%s\n", addr)
	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

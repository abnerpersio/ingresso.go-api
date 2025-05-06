package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"ingresso.go/config"
	"ingresso.go/routes"
)

func main() {
	godotenv.Load()

	port := config.GetEnv("PORT", "8080")
	address := ":" + port

	router := routes.Register()

	fmt.Printf("Server starting on http://localhost%s\n", address)
	err := http.ListenAndServe(address, router)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

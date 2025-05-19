package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"ingresso.go/infra/config"
	"ingresso.go/infra/services"
	"ingresso.go/routes"
)

func initEnv() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}
}

func main() {
	initEnv()

	port := config.GetEnv("PORT", "8080")
	address := ":" + port

	cognitoService := services.NewCognitoService(services.CognitoConfig{
		UserPoolID:      config.GetEnv("COGNITO_USER_POOL_ID"),
		AppClientID:     config.GetEnv("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: config.GetEnv("COGNITO_APP_CLIENT_SECRET"),
		AppPoolDomain:   config.GetEnv("COGNITO_APP_POOL_DOMAIN"),
	})

	router := routes.Register(routes.RouterParams{
		Cognito: cognitoService,
	})

	fmt.Printf("Server starting on http://localhost%s\n", address)
	err := http.ListenAndServe(address, router)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

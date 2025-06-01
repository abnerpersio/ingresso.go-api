package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"ingresso.go/internal/infra/config"
	"ingresso.go/internal/infra/services"
	"ingresso.go/internal/routes"
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

	cognitoService := services.NewCognitoService(services.CognitoConfig{
		UserPoolID:      config.GetEnv("COGNITO_USER_POOL_ID"),
		AppClientID:     config.GetEnv("COGNITO_APP_CLIENT_ID"),
		AppClientSecret: config.GetEnv("COGNITO_APP_CLIENT_SECRET"),
		AppPoolDomain:   config.GetEnv("COGNITO_APP_POOL_DOMAIN"),
	})

	router := routes.Register(routes.RouterParams{Cognito: cognitoService})

	port := ":" + config.GetEnv("PORT", "8080")
	err := router.Run(port)

	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

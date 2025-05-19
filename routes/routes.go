package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	handlers "ingresso.go/domain"
	"ingresso.go/domain/auth"
	"ingresso.go/infra/services"
)

type RouterParams struct {
	Cognito *services.CognitoService
}

func Register(params RouterParams) *mux.Router {
	router := mux.NewRouter()

	authHandler := &auth.AuthHandler{Cognito: params.Cognito}

	router.HandleFunc("/v1/health", handlers.GetHealth).Methods(http.MethodGet)
	router.HandleFunc("/v1/auth/sign-in", authHandler.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/code", authHandler.ExchangeCode).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/refresh-token", authHandler.RefreshToken).Methods(http.MethodPost)

	return router
}

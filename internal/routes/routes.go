package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"ingresso.go/internal/domain"
	"ingresso.go/internal/domain/auth"
	"ingresso.go/internal/infra/services"
)

type RouterParams struct {
	Cognito *services.CognitoService
}

func Register(params RouterParams) *mux.Router {
	router := mux.NewRouter()

	authHandler := &auth.AuthHandler{Cognito: params.Cognito}

	router.HandleFunc("/v1/health", domain.GetHealth).Methods(http.MethodGet)
	router.HandleFunc("/v1/auth/sign-in", authHandler.SignIn).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/code", authHandler.ExchangeCode).Methods(http.MethodPost)
	router.HandleFunc("/v1/auth/refresh-token", authHandler.RefreshToken).Methods(http.MethodPost)

	return router
}

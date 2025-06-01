package routes

import (
	"github.com/gin-gonic/gin"

	"ingresso.go/internal/domain"
	domain_auth "ingresso.go/internal/domain/auth"
	domain_user "ingresso.go/internal/domain/user"
	"ingresso.go/internal/infra/middlewares"
	"ingresso.go/internal/infra/services"
)

type RouterParams struct {
	Cognito *services.CognitoService
}

func Register(params RouterParams) *gin.Engine {
	router := gin.Default()

	authHandler := &domain_auth.AuthHandler{Cognito: params.Cognito}

	router.Use(middlewares.CorsMiddleware())

	router.GET("/v1/health", domain.GetHealth)
	router.POST("/v1/auth/sign-in", authHandler.SignIn)
	router.POST("/v1/auth/code", authHandler.ExchangeCode)
	router.POST("/v1/auth/refresh-token", authHandler.RefreshToken)

	authMiddleware := middlewares.Auth{Cognito: params.Cognito}
	authorized := router.Group("/", authMiddleware.Middleware())
	authorized.GET("/v1/user/profile", domain_user.GetProfile)

	return router
}

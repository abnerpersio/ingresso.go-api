package auth

import (
	"ingresso.go/internal/infra/services"
)

type AuthHandler struct {
	Cognito *services.CognitoService
}

package auth

import (
	"ingresso.go/services"
)

type AuthHandler struct {
	Cognito *services.CognitoService
}

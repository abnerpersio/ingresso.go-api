package auth

import (
	"ingresso.go/infra/services"
)

type AuthHandler struct {
	Cognito *services.CognitoService
}

package domain_auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/services"
)

type RefreshTokenInput struct {
	Email        string `json:"email"`
	RefreshToken string `json:"refresh_token"`
}

func (auth *AuthHandler) RefreshToken(c *gin.Context) {
	var body RefreshTokenInput
	err := c.ShouldBind(&body)

	if err != nil {
		services.SendError(c, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authResp, err := auth.Cognito.RefreshToken(body.Email, body.RefreshToken)

	if err != nil {
		services.SendError(c, "Authentication failed", http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, authResp.AuthenticationResult)
}

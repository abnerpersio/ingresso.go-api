package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/services"
)

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (auth *AuthHandler) SignIn(c *gin.Context) {
	var body SignInInput
	err := c.ShouldBind(&body)

	if err != nil {
		services.SendError(c, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authResp, err := auth.Cognito.AuthenticateUser(body.Username, body.Password)

	if err != nil {
		services.SendError(c, "Authentication failed", http.StatusUnauthorized)
		return
	}

	c.JSON(http.StatusOK, authResp.AuthenticationResult)
}

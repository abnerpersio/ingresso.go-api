package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/services"
)

type RefreshTokenInput struct {
	Username     string `json:"username" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (auth *AuthHandler) RefreshToken(c *gin.Context) {
	var body RefreshTokenInput
	err := c.ShouldBind(&body)

	if err != nil {
		services.SendError(c, "Invalid request payload", http.StatusBadRequest)
		return
	}

	authResp, err := auth.Cognito.RefreshToken(body.Username, body.RefreshToken)

	if err != nil {
		fmt.Println("Error refreshing token:", err)
		services.SendError(c, "Authentication failed", http.StatusUnauthorized)
		return
	}

	fmt.Printf("User resp token %p", authResp)

	c.JSON(http.StatusOK, authResp.AuthenticationResult)
}

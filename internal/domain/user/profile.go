package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/interfaces"
	"ingresso.go/internal/infra/middlewares"
)

func GetProfile(c *gin.Context) {
	user := c.MustGet(middlewares.UserContextKey).(interfaces.User)

	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"id":       user.Id,
			"name":     user.Name,
			"email":    user.Email,
			"provider": user.Provider,
		},
	})
}

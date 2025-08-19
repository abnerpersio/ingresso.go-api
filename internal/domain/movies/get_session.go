package movies

import (
	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/repositories"
)

func GetSession(c *gin.Context) {
	sessionId := c.Param("sessionId")
	repo := &repositories.SessionPGRepository{}

	result, err := repo.Find(sessionId)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": result})
}

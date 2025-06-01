package services

import (
	"github.com/gin-gonic/gin"
)

func SendError(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"message": message})
}

package movies

import (
	"github.com/gin-gonic/gin"
)

func ListSessions(c *gin.Context) {
	// service := services.NewMovieService()
	// movieId := c.Param("movieId")

	// // result := service.ListSessions(movieId)

	// c.JSON(200, gin.H{"data": result})
	c.JSON(200, gin.H{"data": ""})
}

package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/services"
)

func GetMovie(c *gin.Context) {
	service := services.NewMovieService()
	movieId := c.Param("movieId")
	result, err := service.Get(movieId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movie"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

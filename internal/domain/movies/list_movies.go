package movies

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/services"
)

func ListMovies(c *gin.Context) {
	service := services.NewMovieService()
	result, err := service.List()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

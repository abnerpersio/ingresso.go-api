package movies

import (
	"time"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/repositories"
)

func ListSessions(c *gin.Context) {
	movieId := c.Param("movieId")
	date := c.DefaultQuery("date", time.Now().Format(time.DateOnly))
	repo := &repositories.SessionPGRepository{}

	count, err := repo.Count(movieId, date)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if count == 0 {
		repo.Create(repositories.CreateSessionInput{MovieID: movieId, Date: date, StartTime: "18:15", Room: "1"})
		repo.Create(repositories.CreateSessionInput{MovieID: movieId, Date: date, StartTime: "22:30", Room: "1"})
		repo.Create(repositories.CreateSessionInput{MovieID: movieId, Date: date, StartTime: "17:45", Room: "2"})
		repo.Create(repositories.CreateSessionInput{MovieID: movieId, Date: date, StartTime: "20:00", Room: "2"})
	}

	result, err := repo.List(movieId, date)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": result})
}

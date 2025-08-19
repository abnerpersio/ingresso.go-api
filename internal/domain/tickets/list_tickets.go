package tickets

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/interfaces"
	"ingresso.go/internal/infra/middlewares"
	"ingresso.go/internal/infra/repositories"
)

func ListTickets(c *gin.Context) {
	repo := &repositories.TicketPGRepository{}

	user := c.MustGet(middlewares.UserContextKey).(interfaces.User)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	result, err := repo.ListByUser(repositories.ListTicketByUserInput{
		UserID:  user.Id,
		Page:    page,
		PerPage: perPage,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	total, err := repo.CountByUser(user.Id)

	var lastPage int

	if total > 0 {
		lastPage = int((total + perPage - 1) / perPage)
	} else {
		lastPage = 1
	}

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": gin.H{
		"tickets": result,
		"meta": gin.H{
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
			"total":    total,
		},
	}})
}

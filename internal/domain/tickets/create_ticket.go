package tickets

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ingresso.go/internal/infra/interfaces"
	"ingresso.go/internal/infra/middlewares"
	"ingresso.go/internal/infra/repositories"
	"ingresso.go/internal/infra/services"
)

type CreateTicketInput struct {
	SessionID string `json:"session_id" binding:"required"`
	Seats     string `json:"seats" binding:"required"`
}

func CreateTicket(c *gin.Context) {
	repo := &repositories.TicketPGRepository{}

	user := c.MustGet(middlewares.UserContextKey).(interfaces.User)

	var body CreateTicketInput
	err := c.ShouldBind(&body)

	if err != nil {
		services.SendError(c, "Invalid request payload", http.StatusBadRequest)
		return
	}

	result, err := repo.Create(repositories.CreateTicketInput{
		SessionID: body.SessionID,
		UserID:    user.Id,
		Seats:     body.Seats,
		Email:     user.Email,
		// TODO: calculate amount based on seats
		Amount: 1000,
	})

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"data": result})
}

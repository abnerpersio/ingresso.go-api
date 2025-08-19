package repositories

import (
	"context"

	"ingresso.go/internal/infra/config"
)

type Ticket struct {
	ID        string `db:"id" json:"id"`
	SessionID string `db:"session_id" json:"session_id"`
	UserID    string `db:"user_id" json:"user_id"`
	Seats     string `db:"seats" json:"seats"`
	Email     string `db:"email" json:"email"`
	Amount    int    `db:"amount" json:"amount"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type TicketRepository interface {
	Create(input CreateTicketInput) (string, error)
	ListByUser(input ListTicketByUserInput) ([]Ticket, error)
	CountByUser(userId string) (int, error)
}

type TicketPGRepository struct {
}

type CreateTicketInput struct {
	SessionID string
	UserID    string
	Seats     string
	Email     string
	Amount    int
}

type ListTicketByUserInput struct {
	UserID  string
	Page    int
	PerPage int
}

func (repo *TicketPGRepository) Create(input CreateTicketInput) (string, error) {
	db := config.GetDatabase()
	var id string
	err := db.QueryRow(context.Background(), "INSERT INTO ticket (session_id, user_id, seats, email, amount) VALUES ($1, $2, $3, $4, $5) RETURNING id", input.SessionID, input.UserID, input.Seats, input.Email, input.Amount).Scan(&id)

	if err != nil {
		return "", err
	}

	defer db.Close(context.Background())

	return id, nil
}

func (repo *TicketPGRepository) ListByUser(input ListTicketByUserInput) ([]Ticket, error) {
	db := config.GetDatabase()

	page := input.Page | 1
	perPage := input.PerPage | 10
	offset := (page - 1) * perPage
	rows, err := db.Query(context.Background(), "SELECT id, session_id, user_id, seats, email, amount, created_at FROM ticket where user_id = $1 LIMIT $2 OFFSET $3", input.UserID, perPage, offset)

	if err != nil {
		return []Ticket{}, err
	}

	defer rows.Close()
	var result []Ticket

	for rows.Next() {
		var ticket Ticket

		if err := rows.Scan(&ticket.ID, &ticket.SessionID, &ticket.UserID, &ticket.Seats, &ticket.Email, &ticket.Email, &ticket.Amount, &ticket.CreatedAt); err != nil {
			return []Ticket{}, err
		}

		result = append(result, ticket)
	}

	defer db.Close(context.Background())

	return result, nil
}

func (repo *TicketPGRepository) CountByUser(userId string) (int, error) {
	db := config.GetDatabase()

	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM ticket where user_id = $1", userId).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

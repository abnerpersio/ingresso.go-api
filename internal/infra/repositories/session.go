package repositories

import (
	"context"

	"ingresso.go/internal/infra/config"
)

type Session struct {
	ID        string `db:"id" json:"id"`
	MovieID   string `db:"movie_id" json:"movie_id"`
	StartTime string `db:"start_time" json:"start_time"`
	Date      string `db:"date" json:"date"`
	Room      string `db:"room" json:"room"`
	CreatedAt string `db:"created_at" json:"created_at"`
}

type SessionRepository interface {
	Count(movieId string, date string) (int, error)
	List(movieId string, date string) ([]Session, error)
	Find(sessionId string) (Session, error)
	Create(input CreateSessionInput) (string, error)
}

type SessionPGRepository struct {
}

type CreateSessionInput struct {
	MovieID   string
	Date      string
	StartTime string
	Room      string
}

func (repo *SessionPGRepository) List(movieId string, date string) ([]Session, error) {
	db := config.GetDatabase()
	rows, err := db.Query(context.Background(), "SELECT id, movie_id, start_time, date::text, room, created_at::text from session where movie_id = $1 and date = $2", movieId, date)

	if err != nil {
		return []Session{}, err
	}

	defer rows.Close()
	var result []Session

	for rows.Next() {
		var session Session

		if err := rows.Scan(&session.ID, &session.MovieID, &session.StartTime, &session.Date, &session.Room, &session.CreatedAt); err != nil {
			return []Session{}, err
		}

		result = append(result, session)
	}

	defer db.Close(context.Background())

	return result, nil
}

func (repo *SessionPGRepository) Count(movieId string, date string) (int64, error) {
	db := config.GetDatabase()

	var count int64
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) from session where movie_id = $1 and date = $2", movieId, date).Scan(&count)

	if err != nil {
		return 0, err
	}

	defer db.Close(context.Background())

	return count, nil
}

func (repo *SessionPGRepository) Find(sessionId string) (Session, error) {
	db := config.GetDatabase()

	var session Session

	err := db.QueryRow(context.Background(), "SELECT id, movie_id, start_time, date::text, room, created_at::text from session where id = $1", sessionId).Scan(&session.ID, &session.MovieID, &session.StartTime, &session.Date, &session.Room, &session.CreatedAt)

	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (repo *SessionPGRepository) Create(input CreateSessionInput) (string, error) {
	db := config.GetDatabase()
	var id string
	err := db.QueryRow(context.Background(), "INSERT INTO session (movie_id, date, start_time, room) VALUES ($1, $2, $3, $4) RETURNING id", input.MovieID, input.Date, input.StartTime, input.Room).Scan(&id)

	if err != nil {
		return "", err
	}

	defer db.Close(context.Background())

	return id, nil
}

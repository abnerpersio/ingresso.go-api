package config

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetDatabase() *pgx.Conn {
	db, err := pgx.Connect(context.Background(), GetEnv("DATABASE_URL"))

	if err != nil {
		fmt.Printf("unable to connect to database: %v\n", err)
		return nil
	}

	return db
}

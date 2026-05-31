package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Open creates and verifies a PostgreSQL connection using the provided DSN.
func Open(dsn string) (*sql.DB, error) {
	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := database.Ping(); err != nil {
		_ = database.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return database, nil
}

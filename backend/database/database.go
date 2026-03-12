package database

import (
	"database/sql"
	"log/slog"

	_ "github.com/lib/pq"
)

func NewDatabaseConnection(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		slog.Error("[ERROR]", "Failed to find database", err)
	}

	if err = db.Ping(); err != nil {
		slog.Error("[ERROR]", "Failed to establish connection to db", err)
	}

	slog.Info("[DEBUG] Database Established")
	return db
}
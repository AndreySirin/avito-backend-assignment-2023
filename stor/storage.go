package stor

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func NewStorage(log *slog.Logger, user, password, dbname, host, port string) (*Storage, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port,
	)
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &Storage{
		db:  sqlDB,
		log: log,
	}, nil
}

func (s *Storage) Close() error { return s.db.Close() }

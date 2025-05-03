package storage

import (
	"database/sql"
	"fmt"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	lg *logger.MyLogger
	db *sql.DB
}

func New(log *logger.MyLogger, user, password, dbname, host, port string) (*Storage, error) {
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, dbname, host, port,
	)
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return &Storage{
		lg: log,
		db: sqlDB,
	}, nil
}

func (s *Storage) Close() error { return s.db.Close() }

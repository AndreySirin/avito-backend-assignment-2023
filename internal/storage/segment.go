package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
)

type SegmentStorage struct {
	lg *logger.MyLogger
	db *sql.DB
}

func NewSegment(db *Storage) *SegmentStorage {
	return &SegmentStorage{
		lg: db.lg,
		db: db.db,
	}
}

func (s *SegmentStorage) CreateSegment(ctx context.Context, segment entity.Segment) (int, error) {
	err := s.db.QueryRowContext(ctx, "INSERT INTO segments (title) VALUES ($1) RETURNING id_segment", segment.Title).
		Scan(&segment)
	if err != nil {
		return 0, fmt.Errorf("error from CreateSegmen %s", err)
	}
	return segment.ID, nil
}

func (s *SegmentStorage) DeleteSegment(ctx context.Context, segment entity.Segment) error {
	res, err := s.db.ExecContext(
		ctx,
		"DELETE FROM segments WHERE id_segment=$1",
		segment.ID,
	)
	if err != nil {
		return fmt.Errorf("error from DeleteSegg %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	return nil
}

func (s *SegmentStorage) UpDateSegment(ctx context.Context, segment entity.Segment) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %v", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM segment WHERE id_segment=$1)", segment.ID).
		Scan(&exists)
	if err != nil {
		return fmt.Errorf("check segment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("segment id %d does not exist", segment.ID)
	}

	_, err = tx.ExecContext(
		ctx,
		"UPDATE segment SET title=$1 WHERE id_segment=$2",
		segment.Title,
		segment.ID,
	)
	if err != nil {
		return fmt.Errorf("update segment: %w", err)
	}
	return nil
}

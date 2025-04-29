package stor

import (
	"context"
	"fmt"
)

type segment struct {
	id_segment int
	title      string
}

func (s *Storage) CreateSegment(ctx context.Context, segment segment) (int, error) {
	err := s.db.QueryRowContext(ctx, "INSERT INTO segments (title) VALUES ($1) RETURNING id_segment", segment.title).Scan(&segment.id_segment)
	if err != nil {
		return 0, fmt.Errorf("error from CreateSegmen %s", err)
	}
	return segment.id_segment, nil
}
func (s *Storage) DeleteSegment(ctx context.Context, segment segment) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM segments WHERE id_segment=$1", segment.id_segment)
	if err != nil {
		return fmt.Errorf("error from DeleteSegg %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	return nil
}

func (s *Storage) UpDateSegment(ctx context.Context, segment segment) (err error) {

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM segment WHERE id_segment=$1)", segment.id_segment).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check segment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("segment id %d does not exist", segment.id_segment)
	}

	_, err = tx.ExecContext(ctx, "UPDATE segment SET title=$1 WHERE id_segment=$2", segment.title, segment.id_segment)
	if err != nil {
		return fmt.Errorf("update segment: %w", err)
	}
	return nil
}

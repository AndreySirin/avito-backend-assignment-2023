package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"time"
)

func (s *Storage) CreateSegment(ctx context.Context, segment entity.Segment) (uuid.UUID, error) {
	err := segment.Validate()
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %v", ErrNotValid, err)
	}
	query, args, err := sq.Insert("segments").
		Columns(
			"id",
			"title",
			"description",
			"auto_user_prc",
			"create_at",
			"update_at",
			"delete_at",
		).Values(
		segment.ID,
		segment.Title,
		segment.Description,
		segment.AutoUserPrc,
		segment.CreatedAt,
		segment.UpdatedAt,
		segment.DeletedAt,
	).Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()

	err = s.db.QueryRowContext(ctx, query, args...).Scan(&segment.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error from CreateSegment %v", err)
	}
	return segment.ID, nil
}

func (s *Storage) GetSegment(ctx context.Context, id uuid.UUID) (*entity.Segment, error) {
	var segment entity.Segment
	query, args, err := sq.Select(
		"title",
		"description",
		"auto_user_prc",
	).From("segments").
		Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error from GetSegment %v", err)
	}

	err = s.db.QueryRowContext(ctx, query, args...).Scan(segment.Title, segment.Description, segment.AutoUserPrc)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("err not rows %v", err)
		}
	}
	return &segment, nil
}

func (s *Storage) ListSegments(ctx context.Context) ([]entity.Segment, error) {
	var segments []entity.Segment

	query, args, err := sq.Select("*").
		From("segments").
		Where(sq.Eq{"delete_at": nil}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error for build query %v", err)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error for query %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var segment entity.Segment
		err = rows.Scan(&segment.Title, &segment.Description, &segment.AutoUserPrc)
		if err != nil {
			return nil, fmt.Errorf("error from ListSegments %v", err)
		}
		segments = append(segments, segment)

		err = rows.Err()
		if err != nil {
			return nil, fmt.Errorf("error for rows %v", err)
		}
	}
	return segments, nil

}

func (s *Storage) DeleteSegment(ctx context.Context, id uuid.UUID) error {

	t := time.Now()
	query, arg, err := sq.Update("segments").
		Set("delete_at", t).
		Where(sq.Eq{"id": id}).
		Where(sq.Eq{"delete_at": nil}).ToSql()
	res, err := s.db.ExecContext(ctx, query, arg...)

	if err != nil {
		return fmt.Errorf("error from DeleteSegment %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	return nil
}

func (s *Storage) UpDateSegment(ctx context.Context, segment entity.Segment) (err error) {
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
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM segments WHERE id_segment=$1)", segment.ID).
		Scan(&exists)
	if err != nil {
		return fmt.Errorf("check segment existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("segment id %d does not exist", segment.ID)
	}

	query, arg, err := sq.Update("segments").
		Set("title", segment.Title).
		Set("description", segment.Description).
		Set("auto_user_prc", segment.AutoUserPrc).
		Where(sq.Eq{"id": segment.ID}).
		Where(sq.Eq{"delete_at": nil}).ToSql()

	_, err = tx.ExecContext(ctx, query, arg...)
	if err != nil {
		return fmt.Errorf("update segment: %w", err)
	}
	return nil
}

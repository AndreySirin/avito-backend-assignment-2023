package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
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
		).Values(
		segment.ID,
		segment.Title,
		segment.Description,
		segment.AutoUserPrc,
		segment.CreatedAt,
		segment.UpdatedAt,
	).Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return uuid.Nil, fmt.Errorf("error building query: %w", err)
	}

	err = s.Db.QueryRowContext(ctx, query, args...).Scan(&segment.ID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error from CreateSegment %v", err)
	}
	return segment.ID, nil
}

func (s *Storage) GetSegment(ctx context.Context, id uuid.UUID) (*entity.Segment, error) {
	var segment entity.Segment
	query, args, err := sq.Select(
		"id",
		"title",
		"description",
		"auto_user_prc",
		"create_at",
		"update_at",
	).From("segments").Where(sq.And{
		sq.Eq{"id": id},
		sq.Eq{"delete_at": nil}}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error from GetSegment %v", err)
	}

	err = s.Db.QueryRowContext(ctx, query, args...).
		Scan(&segment.ID, &segment.Title, &segment.Description, &segment.AutoUserPrc, &segment.CreatedAt, &segment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("err not rows %v", err)
		}
		return nil, fmt.Errorf("error from GetSegment %v", err)
	}
	return &segment, nil
}

func (s *Storage) ListSegments(ctx context.Context) ([]entity.Segment, error) {
	query, args, err := sq.Select(
		"id",
		"title",
		"description",
		"auto_user_prc",
		"create_at",
		"update_at").
		From("segments").
		Where(sq.Eq{"delete_at": nil}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error for build query %v", err)
	}
	rows, err := s.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error for query %v", err)
	}

	defer func() {
		if errClose := rows.Close(); errClose != nil {
			s.Lg.Error("error closing rows", "error", errClose)
			return
		}
	}()

	var segments []entity.Segment
	for rows.Next() {
		var segment entity.Segment
		err = rows.Scan(
			&segment.ID,
			&segment.Title,
			&segment.Description,
			&segment.AutoUserPrc,
			&segment.CreatedAt,
			&segment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error from ListSegments %v", err)
		}
		segments = append(segments, segment)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error for rows %v", err)
	}
	return segments, nil
}

func (s *Storage) DeleteSegment(ctx context.Context, id uuid.UUID) error {
	query, arg, err := sq.Update("segments").
		Set("delete_at", time.Now()).
		Where(sq.And{
			sq.Eq{"id": id},
			sq.Eq{"delete_at": nil},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error from DeleteSegment %v", err)
	}
	res, err := s.Db.ExecContext(ctx, query, arg...)
	if err != nil {
		return fmt.Errorf("error from DeleteSegment %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error from DeleteSegment %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("no rows affected: segment not found or already deleted")
	}
	return nil
}

func (s *Storage) UpDateSegment(ctx context.Context, segment entity.Segment) (err error) {
	tx, err := s.Db.BeginTx(ctx, nil)
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

	query, _, err := sq.Select("EXISTS(SELECT 1 FROM segments WHERE id = ?)").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error from UpDateSegment %v", err)
	}

	err = tx.QueryRowContext(ctx, query, segment.ID).
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
		Where(sq.And{
			sq.Eq{"id": segment.ID},
			sq.Eq{"delete_at": nil},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("error from UpDateSegment %v", err)
	}

	_, err = tx.ExecContext(ctx, query, arg...)
	if err != nil {
		return fmt.Errorf("update segment: %w", err)
	}
	return nil
}

func (s *Storage) GetIDForSegment(
	ctx context.Context,
	tx *sql.Tx,
	sub []string,
) ([]uuid.UUID, error) {
	var id []uuid.UUID

	query, arg, err := sq.Select("id").
		From("segments").
		Where(sq.And{
			sq.Eq{"delete_at": nil},
			sq.Expr("title = ANY(?)", sub),
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error from GetIDForSegment %v", err)
	}

	rows, err := tx.QueryContext(ctx, query, arg...)
	if err != nil {
		return nil, fmt.Errorf("error selecting segment: %w", err)
	}

	defer func() {
		if errClose := rows.Close(); errClose != nil {
			s.Lg.Error("error closing rows", "error", errClose)
			return
		}
	}()

	for rows.Next() {
		var i uuid.UUID
		err = rows.Scan(&i)
		if err != nil {
			return nil, fmt.Errorf("error scanning segment: %w", err)
		}
		id = append(id, i)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return id, nil
}

func (s *Storage) GetTitleForSegment(ctx context.Context, id []uuid.UUID) (map[uuid.UUID]string, error) {
	title := make(map[uuid.UUID]string)

	query, args, err := sq.Select(
		"id",
		"title").
		From("segments").
		Where(sq.And{
			sq.Eq{"delete_at": nil},
			sq.Expr("id = ANY(?)", id),
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error from GetTitleForSegment %v", err)
	}

	rows, err := s.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error selecting segment: %w", err)
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			s.Lg.Error("error closing rows", "error", errClose)
			return
		}
	}()

	for rows.Next() {
		var i uuid.UUID
		var str string
		err = rows.Scan(&i, &str)
		if err != nil {
			return nil, fmt.Errorf("error scanning segment: %w", err)
		}
		title[i] = str
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", err)
	}
	return title, nil
}

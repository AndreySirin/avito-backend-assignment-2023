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

func (s *Storage) InsertSubscription(
	ctx context.Context,
	tx *sql.Tx,
	subs []entity.Subscription,
) error {
	build := sq.Insert("subscriptions").
		Columns(
			"user_id",
			"segment_id",
			"ttl",
			"is_auto_add",
		)

	for _, sub := range subs {
		build = build.Values(
			sub.UserID,
			sub.SegmentID,
			sub.TTL,
			sub.IsAutoAdded,
		)
	}
	query, args, err := build.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error inserting segment: %w", err)
	}
	return nil
}

func (s *Storage) DeleteSubscription(
	ctx context.Context,
	tx *sql.Tx,
	id uuid.UUID,
	segmentID []uuid.UUID,
) error {
	query, args, err := sq.Delete("subscriptions").
		Where(sq.And{
			sq.Expr("segment_id=ANY(?)", segmentID),
			sq.Eq{"user_id": id},
		}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting segment: %w", err)
	}
	row, err := res.RowsAffected()
	if err != nil || row == 0 {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	return nil
}

func (s *Storage) GetHistorySubscription(
	ctx context.Context,
	date *time.Time,
) ([]entity.HistorySubscription, error) {
	query, args, err := sq.Select(
		"user_id",
		"segment_id",
		"created_at",
		"update_at",
		"delete_at",
	).From("subscriptions").
		Where(sq.Gt{
			"created_at": date,
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error selecting subscriptions: %w", err)
	}
	if date == nil {
		return nil, fmt.Errorf("missing required date parameter")
	}

	rows, err := s.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying subscriptions: %w", err)
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			s.Lg.Error("error closing rows", "error", errClose)
			return
		}
	}()
	var subs []entity.HistorySubscription
	for rows.Next() {
		var sub entity.HistorySubscription
		err = rows.Scan(&sub.UserID, &sub.SegmentID, &sub.CreatedAt, &sub.UpdatedAt, &sub.DeletedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		subs = append(subs, sub)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return subs, nil
}

func (s *Storage) CheckTTLSubscription(ctx context.Context) (int, error) {
	query, args, err := sq.Update("subscriptions").
		Set("delete_at", time.Now()).
		Where(sq.And{
			sq.Eq{"delete_at": nil},
			sq.LtOrEq{"ttl": time.Now()},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return 0, fmt.Errorf("error updating TTL: %w", err)
	}

	res, err := s.Db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("error checking TTL: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error checking TTL: %w", err)
	}
	return int(rows), nil
}

func (s *Storage) GetUserSubscription(ctx context.Context, id uuid.UUID) ([]entity.Subscription, error) {
	var subs []entity.Subscription
	query, args, err := sq.Select(
		"user_id",
		"segment_id",
		"is_auto_add",
	).From("subscriptions").
		Where(sq.And{
			(sq.Eq{"user_id": id}),
			sq.Eq{"delete_at": nil},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error selecting subscription: %w", err)
	}
	rows, err := s.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying subscription: %w", err)
	}
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			s.Lg.Error("error closing rows", "error", errClose)
			return
		}
	}()
	for rows.Next() {
		var u entity.Subscription
		err = rows.Scan(&u.UserID, &u.SegmentID, &u.IsAutoAdded)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		subs = append(subs, u)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return subs, err
}

func (s *Storage) GetUsersIDForSegment(ctx context.Context, segmentID uuid.UUID) ([]uuid.UUID, error) {
	var userIds []uuid.UUID

	query, args, err := sq.Select("user_id").
		From("subscriptions").
		Where(sq.And{
			sq.Eq{"segment_id": segmentID},
			sq.Eq{"is_auto_add": true},
		}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("error selecting users: %w", err)
	}
	rows, err := s.Db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	for rows.Next() {
		var u uuid.UUID
		err = rows.Scan(&u)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		userIds = append(userIds, u)
	}
	return userIds, nil
}

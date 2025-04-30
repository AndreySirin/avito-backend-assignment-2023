package storage

import (
	"context"
	"fmt"
	"time"
)

type Subscription struct {
	IdUser      int       `json:"id_user"`
	NameSegment []string  `json:"name_segment"`
	IdSegment   []int     `json:"id_segment"`
	CreatedAt   time.Time `json:"created_at"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type User_Subscription interface {
	InsertUserInSegment(context.Context, Subscription) (err error)
	DeleteUserInSegment(context.Context, Subscription) (err error)
}

func (s *Storage) InsertUserInSegment(ctx context.Context, subs Subscription) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id_user = $1)", subs.IdUser).
		Scan(&exists)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("user does not exist")
	}

	rows, err := tx.QueryContext(
		ctx,
		"SELECT id_segment FROM segments WHERE title=ANY($1)",
		subs.NameSegment,
	)
	if err != nil {
		return fmt.Errorf("error selecting segment: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			return fmt.Errorf("error scanning segment: %w", err)
		}
		subs.IdSegment = append(subs.IdSegment, i)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate rows: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO subscriptions (id_user,id_segment,expires_at)
SELECT $1, id_segm,$3
FROM UNNEST($2::int[]) AS id_segm`,
		subs.IdUser, subs.IdSegment, subs.ExpiresAt)
	if err != nil {
		return fmt.Errorf("error inserting segment: %w", err)
	}
	return nil
}

func (s *Storage) DeleteUserInSegment(ctx context.Context, subs Subscription) (err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	rows, err := tx.QueryContext(
		ctx,
		"SELECT id_segment FROM segments WHERE title=ANY($1)",
		subs.NameSegment,
	)
	if err != nil {
		return fmt.Errorf("error selecting segment: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			return fmt.Errorf("error scanning segment: %w", err)
		}
		subs.IdSegment = append(subs.IdSegment, i)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("failed to iterate rows: %w", err)
	}
	res, err := tx.ExecContext(
		ctx,
		`DELETE FROM subscriptions WHERE id_user = $1 AND id_segment=ANY($2)`,
		subs.IdUser,
		subs.IdSegment,
	)
	if err != nil {
		return fmt.Errorf("error deleting segment: %w", err)
	}
	row, err := res.RowsAffected()
	if err != nil || row == 0 {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	return nil
}

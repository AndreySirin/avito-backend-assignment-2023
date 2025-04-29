package stor

import (
	"context"
	"fmt"
	"time"
)

type Subscription struct {
	IdSubscription int       `json:"id_subscription"`
	IdUser         int       `json:"id_user"`
	NameSegment    []string  `json:"name_segment"`
	IdSegment      []int     `json:"id_segment"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

type User_Subscription interface {
	InsertUserInSegment(context.Context, Subscription) (id_sabs int, err error)
}

func (s *Storage) InsertUserInSegment(ctx context.Context, subs Subscription) (id_sabs int, err error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	var exists bool
	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", subs.IdUser).Scan(&exists)
	if err != nil {
		return 0, fmt.Errorf("error checking if user exists: %w", err)
	}
	if !exists {
		return 0, fmt.Errorf("user does not exist")
	}

	err = tx.QueryRowContext(ctx,
		`INSERT INTO subscriptions (id_user,expires_at) 
		VALUES ($1,$2) RETURNING id_subscription`, subs.IdUser, subs.ExpiresAt).Scan(&subs.IdSubscription)
	if err != nil {
		return 0, fmt.Errorf("error inserting subscription: %w", err)
	}
	rows, err := tx.QueryContext(ctx, "SELECT id_segment FROM segments WHERE title=ANY($1)", subs.NameSegment)
	if err != nil {
		return 0, fmt.Errorf("error selecting segment: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i int
		err = rows.Scan(&i)
		if err != nil {
			return 0, fmt.Errorf("error scanning segment: %w", err)
		}
		subs.IdSegment = append(subs.IdSegment, i)
	}
	if err = rows.Err(); err != nil {
		return 0, fmt.Errorf("failed to iterate rows: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO subscription_segments(id_subscription,id_segment)
SELECT $1, id_segm
FROM UNNEST($2::int[]) AS id_segm`,
		subs.IdSubscription, subs.IdSegment)
	if err != nil {
		return 0, fmt.Errorf("error inserting segment: %w", err)
	}
	return subs.IdSubscription, nil
}

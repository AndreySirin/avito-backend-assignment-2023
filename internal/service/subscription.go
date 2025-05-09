package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
)

type CreateSubscription struct {
	UserID       uuid.UUID `validate:"required"`
	TitleSegment string    `validate:"required,dive,required"`
	SegmentID    uuid.UUID
	TTL          time.Time
	IsAutoAdded  bool
}

type HistorySubscriptions struct {
	UserID       uuid.UUID  `json:"user_id"       validate:"required"`
	TitleSegment string     `json:"title_segment" validate:"required,dive,required"`
	CreatedAt    time.Time  `json:"created_at"    validate:"required"`
	UpdatedAt    time.Time  `json:"updated_at"    validate:"required"`
	DeletedAt    *time.Time `json:"deleted_at"`
}

func (s *Service) InsertUserInSegments(ctx context.Context, req []CreateSubscription) (err error) {
	tx, err := s.repository.TX(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				s.lg.Error("error to rollback transaction", "err", errRollback)
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				s.lg.Error("error to commit transaction", "err", errCommit)
				err = errCommit
			}
		}
	}()
	Subscription := make([]entity.Subscription, len(req))
	Title := make([]string, len(req))
	for i, r := range req {
		Subscription[i] = entity.Subscription{
			UserID:       r.UserID,
			TitleSegment: r.TitleSegment,
			TTL:          r.TTL,
			IsAutoAdded:  r.IsAutoAdded,
		}
		Title[i] = r.TitleSegment
	}

	UserID := Subscription[0].UserID
	err = s.repository.CheckExistUser(ctx, tx, UserID)
	if err != nil {
		return err
	}

	SegmentID, err := s.repository.GetIDForSegment(ctx, tx, Title)
	if err != nil {
		return err
	}
	if len(SegmentID) == len(Subscription) {
		for i := range Subscription {
			Subscription[i].SegmentID = SegmentID[i]
		}
	}

	err = s.repository.InsertSubscription(ctx, tx, Subscription)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteUserInSegments(ctx context.Context, req []CreateSubscription) (err error) {
	tx, err := s.repository.TX(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				s.lg.Error("error to rollback transaction", "err", errRollback)
			}
		} else {
			errCommit := tx.Commit()
			if errCommit != nil {
				s.lg.Error("error to commit transaction", "err", errCommit)
				err = errCommit
			}
		}
	}()

	Title := make([]string, len(req))
	for i, r := range req {
		Title[i] = r.TitleSegment
	}
	userID := req[0].UserID

	segmentID, err := s.repository.GetIDForSegment(ctx, tx, Title)
	if err != nil {
		return fmt.Errorf("error for get id segment :%v", err)
	}
	err = s.repository.DeleteSubscription(ctx, tx, userID, segmentID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetHistorySubscriptions(
	ctx context.Context,
	date *time.Time,
) (res []HistorySubscriptions, err error) {
	req, err := s.repository.GetHistorySubscription(ctx, date)
	if err != nil {
		return nil, err
	}
	SegmentID := make([]uuid.UUID, len(req))
	for i := 0; i < len(req); i++ {
		SegmentID[i] = req[i].SegmentID
	}

	SegmentTitle, err := s.repository.GetTitleForSegment(ctx, SegmentID)
	if err != nil {
		return nil, err
	}
	history := make([]HistorySubscriptions, len(req))
	for i := 0; i < len(req); i++ {
		cr := req[i].CreatedAt.Truncate(time.Second)

		history[i] = HistorySubscriptions{
			UserID:       req[i].UserID,
			TitleSegment: SegmentTitle[SegmentID[i]],
			CreatedAt:    cr,
		}
		if req[i].DeletedAt != nil {
			dl := req[i].DeletedAt.Truncate(time.Second)
			history[i].DeletedAt = &dl
		}
	}
	return history, nil
}

func (s *Service) CheckTTLSubscriptions(ctx context.Context) (int, error) {
	rows, err := s.repository.CheckTTLSubscription(ctx)
	if err != nil {
		return -1, err
	}
	return rows, nil
}

func (s *Service) GetUsersSubscription(
	ctx context.Context,
	userID uuid.UUID,
) ([]entity.Subscription, error) {
	sub, err := s.repository.GetUserSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}
	id := make([]uuid.UUID, len(sub))
	for i, r := range sub {
		id[i] = r.SegmentID
	}
	title, err := s.repository.GetTitleForSegment(ctx, id)
	if err != nil {
		return nil, err
	}
	for i, r := range sub {
		sub[i].TitleSegment = title[r.SegmentID]
	}
	return sub, nil
}

func (s *Service) GetUsersIDsForSegment(
	ctx context.Context,
	segmentID uuid.UUID,
) ([]uuid.UUID, error) {
	return s.repository.GetUsersIDForSegment(ctx, segmentID)
}

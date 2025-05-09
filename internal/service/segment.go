package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/validator"
)

type CreateSegment struct {
	Title       string `validate:"required"`
	Description string `validate:"required"`
	AutoUserPrc uint8  `validate:"gte=0,lte=100"`
}

func (c *CreateSegment) Validate() error { return validator.Validator.Struct(c) }

type UpdateSegmentRequest struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required"`
	Description string    `validate:"required"`
	AutoUserPrc uint8
}

func (u *UpdateSegmentRequest) Validate() error { return validator.Validator.Struct(u) }

func (s *Service) CreateSegment(
	ctx context.Context,
	request CreateSegment,
) (uuid.UUID, error) {
	err := request.Validate()
	if err != nil {
		return uuid.Nil, err
	}

	tx, err := s.repository.TX(ctx)
	if err != nil {
		return uuid.Nil, err
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

	segment := entity.Segment{
		ID:          uuid.New(),
		Title:       request.Title,
		Description: request.Description,
		AutoUserPrc: request.AutoUserPrc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	id, err := s.repository.CreateSegment(ctx, segment)
	if err != nil {
		return uuid.Nil, err
	}
	if segment.AutoUserPrc == 0 {
		return id, nil
	}

	allUsersID, err := s.repository.ListUsersID(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	AUP := float32(segment.AutoUserPrc) * 0.01
	quantity := int(math.Round(float64(len(allUsersID)) * float64(AUP)))
	sliceUserID := allUsersID[0:quantity]
	subscription := make([]entity.Subscription, len(sliceUserID))
	for i, iD := range sliceUserID {
		subscription[i] = entity.Subscription{
			UserID:      iD,
			SegmentID:   segment.ID,
			IsAutoAdded: true,
		}
	}
	err = s.repository.InsertSubscription(ctx, tx, subscription)
	if err != nil {
		return uuid.Nil, err
	}
	return segment.ID, nil
}

func (s *Service) GetSegment(ctx context.Context, id uuid.UUID) (*entity.Segment, error) {
	segment, err := s.repository.GetSegment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error servis %v", err)
	}
	return segment, nil
}

func (s *Service) ListSegment(ctx context.Context) ([]entity.Segment, error) {
	segments, err := s.repository.ListSegments(ctx)
	if err != nil {
		return nil, fmt.Errorf("error servis %v", err)
	}
	return segments, nil
}

func (s *Service) UpdateSegment(ctx context.Context, request UpdateSegmentRequest) (err error) {
	err = request.Validate()
	if err != nil {
		return fmt.Errorf("%w,%v", ErrNotValid, err)
	}
	segment := entity.Segment{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
		AutoUserPrc: request.AutoUserPrc,
	}
	err = s.repository.UpDateSegment(ctx, segment)
	if err != nil {
		return fmt.Errorf("error servis %v", err)
	}
	return nil
}

func (s *Service) DeleteSegment(ctx context.Context, id uuid.UUID) error {
	err := s.repository.DeleteSegment(ctx, id)
	if err != nil {
		return fmt.Errorf("error DeleteSegment %v", err)
	}
	return nil
}

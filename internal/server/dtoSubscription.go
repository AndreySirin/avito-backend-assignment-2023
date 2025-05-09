package server

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/validator"
)

type SubscriptionRequest struct {
	UserID       uuid.UUID `json:"id_user"      validate:"required"`
	TitleSegment []string  `json:"name_segment" validate:"required,dive,required"`
	TTL          []string  `json:"ttl"`
	IsAutoAdded  []bool    `json:"auto_added"`
}

type HistorySubscription struct {
	Time string `json:"time" validate:"required"`
}

func (s *SubscriptionRequest) Validate() error {
	if s.UserID == uuid.Nil {
		return errors.New("id_user is required")
	}
	return validator.Validator.Struct(s)
}

func (s *SubscriptionRequest) ToService() ([]service.CreateSubscription, error) {
	var err error
	var SliceSubs []service.CreateSubscription
	for i := 0; i < len(s.TitleSegment); i++ {
		var createSubs service.CreateSubscription
		createSubs.UserID = s.UserID
		createSubs.TitleSegment = s.TitleSegment[i]
		createSubs.TTL, err = time.Parse("2006-01-02 15:04:05", s.TTL[i])
		if err != nil {
			return nil, fmt.Errorf("invalid TTL at index %d: %w", i, err)
		}
		createSubs.IsAutoAdded = s.IsAutoAdded[i]
		SliceSubs = append(SliceSubs, createSubs)
	}
	return SliceSubs, nil
}

func (h *HistorySubscription) ToService() (*time.Time, error) {
	date, err := time.Parse("2006-01-02 15:04:05", h.Time)
	if err != nil {
		return nil, fmt.Errorf("invalid parse time: %v", err)
	}
	return &date, nil
}

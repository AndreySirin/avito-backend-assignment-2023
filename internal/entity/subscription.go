package entity

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	UserID       uuid.UUID `validate:"required"`
	SegmentID    uuid.UUID `validate:"required"`
	TitleSegment string    `validate:"required,required"`
	TTL          time.Time
	IsAutoAdded  bool
	CreatedAt    time.Time `validate:"required"`
	UpdatedAt    time.Time `validate:"required"`
	DeletedAt    *time.Time
}

type HistorySubscription struct {
	UserID    uuid.UUID `validate:"required"`
	SegmentID uuid.UUID `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
	DeletedAt *time.Time
}

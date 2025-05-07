package entity

import (
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/validator"
	"github.com/google/uuid"
	"time"
)

type Segment struct {
	ID          uuid.UUID `validate:"required"`
	Title       string    `validate:"required"`
	Description string    `validate:"required"`
	AutoUserPrc uint8     `validate:"required,gte=0,lte=100"`
	CreatedAt   time.Time `validate:"required"`
	UpdatedAt   time.Time `validate:"required"`
	DeletedAt   time.Time `validate:"omitempty"`
}

func (s *Segment) Validate() error { return validator.Validator.Struct(s) }

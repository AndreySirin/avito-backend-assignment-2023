package server

import (
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/validator"
	"github.com/google/uuid"
)

type createSegmentRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	AutoUserPrc uint8  `json:"auto_user_prc" validate:"required,gte=0,lte=100"`
}

type updateSegmentRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	AutoUserPrc uint8     `json:"auto_user_prc" validate:"required,gte=0,lte=100"`
}

func (c *createSegmentRequest) Validate() error { return validator.Validator.Struct(c) }

func (u *updateSegmentRequest) Validate() error { return validator.Validator.Struct(u) }

func (c *createSegmentRequest) ToService(request createSegmentRequest) service.CreateSegmentRequest {
	return service.CreateSegmentRequest{
		Title:       request.Title,
		Description: request.Description,
		AutoUserPrc: request.AutoUserPrc,
	}
}
func (c *updateSegmentRequest) ToService(request updateSegmentRequest) service.UpdateSegmentRequest {
	return service.UpdateSegmentRequest{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
		AutoUserPrc: request.AutoUserPrc,
	}
}

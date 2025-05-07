package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/validator"
)

type createUserRequest struct {
	FullName    string `json:"fullName"    validate:"required"`
	Gender      string `json:"gender"      validate:"required,oneof=male female"`
	DateOfBirth string `json:"dateOfBirth" validate:"required"`
}

type updateUserRequest struct {
	Id          uuid.UUID `json:"id"    validate:"required"`
	FullName    string    `json:"fullName"    validate:"required"`
	Gender      string    `json:"gender"      validate:"required,oneof=male female"`
	DateOfBirth string    `json:"dateOfBirth" validate:"required"`
}

func (r *createUserRequest) valid() error {
	return validator.Validator.Struct(r)
}

func (r *updateUserRequest) valid() error { return validator.Validator.Struct(r) }

var ErrNotValidDate = errors.New("not valid date")

func (r *createUserRequest) toService() (service.CreateUserRequest, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, r.DateOfBirth)
	if err != nil {
		return service.CreateUserRequest{}, fmt.Errorf("%w: %v", ErrNotValidDate, err)
	}

	reqToService := service.CreateUserRequest{
		FullName:    r.FullName,
		Gender:      r.Gender,
		DateOfBirth: dateOfBirth,
	}

	return reqToService, nil
}

func (u *updateUserRequest) toService() (service.UpdateUserRequest, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, u.DateOfBirth)
	if err != nil {
		return service.UpdateUserRequest{}, fmt.Errorf("%w: %v", ErrNotValidDate, err)
	}
	reqToService := service.UpdateUserRequest{
		Id:          u.Id,
		FullName:    u.FullName,
		Gender:      u.Gender,
		DateOfBirth: dateOfBirth,
	}
	return reqToService, nil

}

package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
)

const module = "service"

type repository interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user entity.User) (err error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	ListUsers(ctx context.Context) ([]entity.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	//
	CreateSegment(ctx context.Context, segment entity.Segment) (uuid.UUID, error)
	UpDateSegment(ctx context.Context, segment entity.Segment) (err error)
	GetSegment(ctx context.Context, id uuid.UUID) (*entity.Segment, error)

	DeleteSegment(ctx context.Context, id uuid.UUID) error

	//

}

type Service struct {
	lg         *slog.Logger
	repository repository
}

func New(lg *slog.Logger, repository repository) *Service {
	return &Service{
		lg:         lg.With("module", module),
		repository: repository,
	}
}

package service

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
)

const module = "service"

type Repository interface {
	CreateUser(ctx context.Context, user entity.User) (uuid.UUID, error)
	UpdateUser(ctx context.Context, user entity.User) (err error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	ListUsers(ctx context.Context) ([]entity.User, error)
	ListUsersID(ctx context.Context) ([]uuid.UUID, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserSubscription(ctx context.Context, id uuid.UUID) ([]entity.Subscription, error)
	CheckExistUser(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
	//
	CreateSegment(ctx context.Context, segment entity.Segment) (uuid.UUID, error)
	UpDateSegment(ctx context.Context, segment entity.Segment) (err error)
	GetSegment(ctx context.Context, id uuid.UUID) (*entity.Segment, error)
	ListSegments(ctx context.Context) ([]entity.Segment, error)
	DeleteSegment(ctx context.Context, id uuid.UUID) error
	GetIDForSegment(ctx context.Context, tx *sql.Tx, sub []string) ([]uuid.UUID, error)
	GetTitleForSegment(ctx context.Context, id []uuid.UUID) (map[uuid.UUID]string, error)
	//
	TX(ctx context.Context) (*sql.Tx, error)
	InsertSubscription(ctx context.Context, tx *sql.Tx, sub []entity.Subscription) (err error)
	DeleteSubscription(ctx context.Context, tx *sql.Tx, id uuid.UUID, SegmentID []uuid.UUID) (err error)
	//
	GetHistorySubscription(ctx context.Context, data *time.Time) ([]entity.HistorySubscription, error)
	CheckTTLSubscription(ctx context.Context) (int, error)
	GetUsersIDForSegment(ctx context.Context, segmentID uuid.UUID) ([]uuid.UUID, error)
}

type Service struct {
	lg         *slog.Logger
	repository Repository
}

func New(lg *slog.Logger, repository Repository) *Service {
	return &Service{
		lg:         lg.With("module", module),
		repository: repository,
	}
}

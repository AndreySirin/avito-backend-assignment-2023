package server

import (
	"context"
	"encoding/json"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
	"net/http"
)

type UserHandle interface {
	CreateUsers(context.Context, entity.User) (int, error)
	UpdateUsers(context.Context, entity.User) (err error)
	DeleteUsers(context.Context, entity.User) error
}

type SegmentHandle interface {
	CreateSegments(context.Context, entity.Segment) (int, error)
	DeleteSegments(context.Context, entity.Segment) error
	UpDateSegments(context.Context, entity.Segment) (err error)
}

type SubscriptionHandle interface {
	InsertUserInSegments(ctx context.Context, subs entity.CreateSubscription) (err error)
	DeleteUserInSegments(ctx context.Context, subs entity.CreateSubscription) (err error)
}

type HNDL struct {
	UserHandler         UserHandle
	SegmentHandler      SegmentHandle
	SubscriptionHandler SubscriptionHandle
	lg                  logger.MyloggerInterface
}

func NewHNDL(lg *logger.MyLogger, userHandler UserHandle, segmentHandler SegmentHandle, subscriptionHandler SubscriptionHandle) *HNDL {
	return &HNDL{
		lg:                  lg,
		UserHandler:         userHandler,
		SegmentHandler:      segmentHandler,
		SubscriptionHandler: subscriptionHandler,
	}
}

func (H *HNDL) UserAddSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var subs entity.CreateSubscription
	err := json.NewDecoder(r.Body).Decode(&subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = H.SubscriptionHandler.InsertUserInSegments(context.Background(), subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode("successfully"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (H *HNDL) UserDeleteSegment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	var subs entity.CreateSubscription
	err := json.NewDecoder(r.Body).Decode(&subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = H.SubscriptionHandler.DeleteUserInSegments(context.Background(), subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode("successfully"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

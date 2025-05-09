package integ

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"

	"github.com/AndreySirin/avito-backend-assignment-2023/internal/server"
)

func (s *SegmentIntegritionTest) TestCreateSubscription() {
	req := server.SubscriptionRequest{
		UserID:       uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		TitleSegment: []string{"Segment C"},
		TTL:          []string{"2025-06-01 10:00:00"},
		IsAutoAdded:  []bool{false},
	}
	body, err := json.Marshal(req)
	s.Suite.NoError(err)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/subscription", bytes.NewReader(body))
	s.handler.HttpServer.Handler.ServeHTTP(recorder, request)
	s.Equal(http.StatusCreated, recorder.Code)
}

func (s *SegmentIntegritionTest) TestDeleteSubscription() {
	req := server.SubscriptionRequest{
		TitleSegment: []string{"Segment A"},
		TTL:          []string{"2025-05-01 12:00:00"},
		IsAutoAdded:  []bool{true},
	}

	body, err := json.Marshal(req)
	s.Suite.NoError(err)

	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	uru := fmt.Sprintf("/api/v1/subscription/%v", id)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodDelete, uru, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	s.handler.HttpServer.Handler.ServeHTTP(recorder, request)
	s.Equal(http.StatusOK, recorder.Code)
}

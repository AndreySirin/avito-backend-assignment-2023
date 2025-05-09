package unit

import (
	"bytes"
	"encoding/json"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/server"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
	"github.com/AndreySirin/avito-backend-assignment-2023/test/unit/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHundlerCreateUser(t *testing.T) {

	type mockBehavior func(*mocks.Srv, service.CreateUserRequest)

	tests := []struct {
		TitleTest          string
		mockBehavior       mockBehavior
		User               server.CreateUserRequest
		user_servis        service.CreateUserRequest
		expectedStatusCode int
	}{
		{
			TitleTest: "good",
			mockBehavior: func(srv *mocks.Srv, request service.CreateUserRequest) {
				srv.On("CreateUser", mock.Anything, request).Return(uuid.New(), nil)
			},
			User: server.CreateUserRequest{
				FullName:    "andrey",
				Gender:      "male",
				DateOfBirth: "1996-02-01",
			},
			user_servis: service.CreateUserRequest{
				FullName:    "andrey",
				Gender:      "male",
				DateOfBirth: time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			TitleTest:    "bad",
			mockBehavior: func(srv *mocks.Srv, request service.CreateUserRequest) {},
			User: server.CreateUserRequest{
				FullName:    "andrey",
				Gender:      "e",
				DateOfBirth: "1996-02-01",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(
			test.TitleTest,
			func(t *testing.T) {
				srv := mocks.NewSrv(t)
				s := server.Server{Service: srv}

				test.mockBehavior(srv, test.user_servis)

				r := chi.NewRouter()
				r.Post("/users", s.HandleCreateUser)

				body, err := json.Marshal(test.User)
				if err != nil {
					return
				}

				recorder := httptest.NewRecorder()
				request := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
				r.ServeHTTP(recorder, request)

				require.Equal(t, test.expectedStatusCode, recorder.Code)
				srv.AssertExpectations(t)
			},
		)
	}
}

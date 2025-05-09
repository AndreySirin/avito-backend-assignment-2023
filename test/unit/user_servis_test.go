package unit

import (
	"context"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/service"
	"github.com/AndreySirin/avito-backend-assignment-2023/test/unit/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestServisCreateUser(t *testing.T) {

	type mockBehavior func(*mocks.Repository)

	tests := []struct {
		name         string
		mockBehavior mockBehavior
		userRequest  service.CreateUserRequest
		wantErr      bool
	}{
		{
			name:    "good request",
			wantErr: false,
			userRequest: service.CreateUserRequest{
				FullName:    "andrey",
				Gender:      "male",
				DateOfBirth: time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(repository *mocks.Repository) {
				repository.On("CreateUser",
					mock.Anything,
					mock.MatchedBy(func(u entity.User) bool {
						return u.FullName == "andrey" &&
							u.Gender == "male" &&
							u.DateOfBirth.Equal(time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC))
					}),
				).Return(uuid.New(), nil)
			},
		},
		{
			name:    "bad gender request",
			wantErr: true,
			userRequest: service.CreateUserRequest{
				FullName:    "andrey",
				Gender:      "",
				DateOfBirth: time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC),
			},
			mockBehavior: func(repository *mocks.Repository) {},
		},
		{
			name:    "bad DateOfBirth request",
			wantErr: true,
			userRequest: service.CreateUserRequest{
				FullName: "andrey",
				Gender:   "male",
			},
			mockBehavior: func(repository *mocks.Repository) {},
		},
		{
			name:         "emty body request",
			wantErr:      true,
			userRequest:  service.CreateUserRequest{},
			mockBehavior: func(repository *mocks.Repository) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repa := mocks.NewRepository(t)
			tt.mockBehavior(repa)
			servis := service.New(logger.New(true), repa)
			id, err := servis.CreateUser(context.Background(), tt.userRequest)
			if tt.wantErr {
				require.Equal(t, uuid.Nil, id)
				require.Error(t, err)
			} else {
				require.NotNil(t, id)
				require.NoError(t, err)
			}
			repa.AssertExpectations(t)
		})
	}
}

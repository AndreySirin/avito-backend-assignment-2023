package unit

import (
	"context"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/entity"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/logger"
	"github.com/AndreySirin/avito-backend-assignment-2023/internal/storage"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestStorage_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	lg := logger.New(true)

	stor := storage.Storage{
		Lg: lg,
		Db: db,
	}

	type mockBehavior func(mock sqlmock.Sqlmock, user entity.User)

	tests := []struct {
		name    string
		wantErr bool
		user    entity.User
		mockb   mockBehavior
	}{
		{
			name:    "good body request",
			wantErr: false,
			user: entity.User{
				ID:          uuid.New(),
				FullName:    "andrey",
				Gender:      "male",
				DateOfBirth: time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   nil,
			},
			mockb: func(mock sqlmock.Sqlmock, user entity.User) {
				mock.ExpectQuery("INSERT INTO users").WithArgs(user.ID, user.FullName, user.Gender, user.DateOfBirth, user.CreatedAt, user.UpdatedAt, user.DeletedAt).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(user.ID))
			},
		},
		{
			name:    "bad gender",
			wantErr: true,
			user: entity.User{
				ID:          uuid.New(),
				FullName:    "andrey",
				Gender:      "ma",
				DateOfBirth: time.Date(1996, time.February, 1, 0, 0, 0, 0, time.UTC),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   nil,
			},
			mockb: func(mock sqlmock.Sqlmock, user entity.User) {},
		},
		{
			name:    "bad data of birth",
			wantErr: true,
			user: entity.User{
				ID:        uuid.New(),
				FullName:  "andrey",
				Gender:    "male",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				DeletedAt: nil,
			},
			mockb: func(mock sqlmock.Sqlmock, user entity.User) {},
		},

		{
			name:    "emty body request",
			wantErr: true,
			user:    entity.User{},
			mockb:   func(mock sqlmock.Sqlmock, user entity.User) {},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockb(mock, tt.user)
			id, err := stor.CreateUser(context.Background(), tt.user)
			if tt.wantErr {
				require.Error(t, err)
				require.Equal(t, uuid.Nil, id)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.user.ID, id)
			}
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

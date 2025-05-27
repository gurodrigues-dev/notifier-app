package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/entity"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetNotificationUsecase_GetNotification(t *testing.T) {
	type args struct {
		id string
	}
	type fields struct {
		notificationError *entity.NotificationError
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *GetNotificationUsecase
		wantErr bool
	}{
		{
			name: "there is to return success",
			args: args{
				id: "1",
			},
			setup: func(t *testing.T) *GetNotificationUsecase {
				notificationError := &entity.NotificationError{
					ID:    1,
					UUID:  "5d91f580-030d-4c5d-b3a6-8383f5829bd3",
					Body:  []byte(`{"message": "error"}`),
					Error: "",
				}
				mock := mocks.NewNotificationRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetNotificationByID", "1").Return(notificationError, nil)
				return NewGetNotificationUsecase(
					mock,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is to return db error",
			args: args{
				id: "1",
			},
			setup: func(t *testing.T) *GetNotificationUsecase {
				mock := mocks.NewNotificationRepository(t)
				logger := mocks.NewLogger(t)
				mock.On("GetNotificationByID", "1").Return(nil, errors.New("db error"))
				return NewGetNotificationUsecase(
					mock,
					logger,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup(t)
			_, err := usecase.GetNotification(tt.args.id)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

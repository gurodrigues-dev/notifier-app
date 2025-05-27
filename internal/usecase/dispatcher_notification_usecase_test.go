package usecase

import (
	"errors"
	"testing"

	"github.com/gurodrigues-dev/notifier-app/internal/infra/contracts"
	"github.com/gurodrigues-dev/notifier-app/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDispatcherUsecase_Execute(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		setup   func(t *testing.T) *DispatcherUsecase
		wantErr bool
	}{
		{
			name: "there return is to success",
			args: args{
				message: `{
					"id": 123,
					"uuid": "550e8400-e29b-41d4-a716-446655440000",
					"title": "Order Confirmation",
					"message": "Your order #12345 has been confirmed.",
					"channels": {
						"1": {
							"id": 1,
							"platform": "email",
							"target_id": "user@example.com",
							"group": "customers"
						},
						"2": {
							"id": 2,
							"platform": "slack",
							"target_id": "webhook_url",
							"group": "customers"
						}
					},
					"event": {
						"name": "OrderPlaced",
						"currency": "BRL",
						"requester": "system",
						"receiver": "user",
						"category": "ecommerce",
						"timestamp": 1716720000,
						"cost_cents": 5000
					},
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)
				ses.On("SendEmail", mock.Anything).Return(nil)
				webhook.On("Post", mock.Anything, mock.Anything, mock.Anything).
					Return(&contracts.HTTPResponse{
						StatusCode: 200,
						Close: func() error {
							return nil
						},
					}, nil)
				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there return is serialize error",
			args: args{
				message: "",
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: true,
		},
		{
			name: "there return is retries greater than 3",
			args: args{
				message: `{
					"retries": 4
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				repository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "when platform not have accessbility",
			args: args{
				message: `{
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				repository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is return to webhook status 500",
			args: args{
				message: `{
					"id": 123,
					"uuid": "550e8400-e29b-41d4-a716-446655440000",
					"title": "Order Confirmation",
					"message": "Your order #12345 has been confirmed.",
					"channels": {
						"2": {
							"id": 2,
							"platform": "slack",
							"target_id": "webhook_url",
							"group": "customers"
						}
					},
					"event": {
						"name": "OrderPlaced",
						"currency": "BRL",
						"requester": "system",
						"receiver": "user",
						"category": "ecommerce",
						"timestamp": 1716720000,
						"cost_cents": 5000
					},
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				webhook.On("Post", mock.Anything, mock.Anything, mock.Anything).
					Return(&contracts.HTTPResponse{
						StatusCode: 500,
						Close: func() error {
							return nil
						},
					}, nil)
				repository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is return to webhook slack error",
			args: args{
				message: `{
					"id": 123,
					"uuid": "550e8400-e29b-41d4-a716-446655440000",
					"title": "Order Confirmation",
					"message": "Your order #12345 has been confirmed.",
					"channels": {
						"2": {
							"id": 2,
							"platform": "slack",
							"target_id": "webhook_url",
							"group": "customers"
						}
					},
					"event": {
						"name": "OrderPlaced",
						"currency": "BRL",
						"requester": "system",
						"receiver": "user",
						"category": "ecommerce",
						"timestamp": 1716720000,
						"cost_cents": 5000
					},
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				webhook.On("Post", mock.Anything, mock.Anything, mock.Anything).
					Return(&contracts.HTTPResponse{
						StatusCode: 500,
						Close: func() error {
							return nil
						},
					}, errors.New("webhook error"))
				repository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is return to webhook discrd error",
			args: args{
				message: `{
					"id": 123,
					"uuid": "550e8400-e29b-41d4-a716-446655440000",
					"title": "Order Confirmation",
					"message": "Your order #12345 has been confirmed.",
					"channels": {
						"2": {
							"id": 2,
							"platform": "discord",
							"target_id": "webhook_url",
							"group": "customers"
						}
					},
					"event": {
						"name": "OrderPlaced",
						"currency": "BRL",
						"requester": "system",
						"receiver": "user",
						"category": "ecommerce",
						"timestamp": 1716720000,
						"cost_cents": 5000
					},
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				webhook.On("Post", mock.Anything, mock.Anything, mock.Anything).
					Return(&contracts.HTTPResponse{
						StatusCode: 500,
						Close: func() error {
							return nil
						},
					}, errors.New("webhook error"))
				repository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "there is return to email error",
			args: args{
				message: `{
					"id": 123,
					"uuid": "550e8400-e29b-41d4-a716-446655440000",
					"title": "Order Confirmation",
					"message": "Your order #12345 has been confirmed.",
					"channels": {
						"1": {
							"id": 1,
							"platform": "email",
							"target_id": "webhook@gmail.com",
							"group": "customers"
						}
					},
					"event": {
						"name": "OrderPlaced",
						"currency": "BRL",
						"requester": "system",
						"receiver": "user",
						"category": "ecommerce",
						"timestamp": 1716720000,
						"cost_cents": 5000
					},
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				ses.On("SendEmail", mock.Anything).Return(errors.New("email error"))
				repository.On("CreateNotification", mock.Anything).Maybe().Return(nil)

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: false,
		},
		{
			name: "when final return to database fails",
			args: args{
				message: `{
					"id": 123,
					"uuid": "550e8400-e29b-41d4-a716-446655440000",
					"title": "Order Confirmation",
					"message": "Your order #12345 has been confirmed.",
					"channels": {
						"1": {
							"id": 1,
							"platform": "email",
							"target_id": "webhook@gmail.com",
							"group": "customers"
						}
					},
					"event": {
						"name": "OrderPlaced",
						"currency": "BRL",
						"requester": "system",
						"receiver": "user",
						"category": "ecommerce",
						"timestamp": 1716720000,
						"cost_cents": 5000
					},
					"retries": 0
				}`,
			},
			setup: func(t *testing.T) *DispatcherUsecase {
				repository := mocks.NewNotificationRepository(t)
				ses := mocks.NewSESIface(t)
				webhook := mocks.NewWebhook(t)
				logger := mocks.NewLogger(t)

				ses.On("SendEmail", mock.Anything).Return(errors.New("email error"))
				repository.On("CreateNotification", mock.Anything).Return(errors.New("db error"))
				logger.On("Errorf", mock.Anything, mock.Anything).Return()

				return NewDispatcherUsecase(
					repository,
					ses,
					webhook,
					logger,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usecase := tt.setup(t)
			err := usecase.Execute(tt.args.message)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

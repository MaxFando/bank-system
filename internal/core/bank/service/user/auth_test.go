package user_test

import (
	"context"
	"github.com/MaxFando/bank-system/internal/core/bank/service/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestNewAuthService(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := user.NewMockAuthRepository(ctrl)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	authService := user.NewAuthService(logger, mockRepo)

	if authService == nil {
		t.Fatal("expected authService to be not nil")
	}
}

func TestAuthService_Register(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name      string
		email     string
		password  string
		setupMock func(m *user.MockAuthRepository)
		want      int32
		err       error
	}{
		{
			name:     "error registering user",
			email:    "email",
			password: "password",
			setupMock: func(m *user.MockAuthRepository) {
				m.EXPECT().
					Register(gomock.Any(), "email", "password").
					Return(int32(0), assert.AnError)
			},
			want: 0,
			err:  assert.AnError,
		},
		{
			name:     "successful registration",
			email:    "email",
			password: "password",
			setupMock: func(m *user.MockAuthRepository) {
				m.EXPECT().
					Register(gomock.Any(), "email", "password").
					Return(int32(1), nil)
			},
			want: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := user.NewMockAuthRepository(ctrl)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			tc.setupMock(mockRepo)

			authService := user.NewAuthService(logger, mockRepo)

			got, err := authService.Register(context.TODO(), tc.email, tc.password)

			if tc.err != nil {
				assert.ErrorAs(t, err, &tc.err)
				assert.Equal(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	testCases := []struct {
		name      string
		email     string
		password  string
		setupMock func(m *user.MockAuthRepository)
		want      string
		err       error
	}{
		{
			name:     "error finding user by email",
			email:    "email",
			password: "password",
			setupMock: func(m *user.MockAuthRepository) {
				m.EXPECT().
					FindByEmail(gomock.Any(), "email").
					Return(nil, assert.AnError)
			},
			want: "",
			err:  assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := user.NewMockAuthRepository(ctrl)
			logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

			tc.setupMock(mockRepo)

			authService := user.NewAuthService(logger, mockRepo)

			got, err := authService.Login(context.TODO(), tc.email, tc.password)

			if tc.err != nil {
				assert.ErrorAs(t, err, &tc.err)
				assert.Equal(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

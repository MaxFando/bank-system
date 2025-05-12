package user

import (
	"context"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

type NullWriter struct{}

func (NullWriter) Write([]byte) (int, error) { return 0, nil }

func TestService_UpdateProfile(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := slog.New(slog.NewTextHandler(NullWriter{}, nil))

	testCases := []struct {
		name     string
		user     *entity.User
		mockFunc func(*MockRepository)
		expected *entity.User
		err      error
	}{
		{
			name: "successful update",
			user: &entity.User{
				ID:       1,
				LastName: "Doe",
			},
			mockFunc: func(m *MockRepository) {
				m.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(&entity.User{
					ID:       1,
					LastName: "Doe",
				}, nil)
			},
		},
		{
			name: "error on update",
			user: &entity.User{
				ID:       1,
				LastName: "Doe",
			},
			mockFunc: func(m *MockRepository) {
				m.EXPECT().UpdateProfile(gomock.Any(), gomock.Any()).Return(nil, assert.AnError)
			},
			err: assert.AnError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewMockRepository(ctrl)
			tc.mockFunc(repo)

			service := NewService(logger, repo)

			got, err := service.UpdateProfile(context.Background(), tc.user)

			if tc.err != nil {
				assert.ErrorAs(t, err, &tc.err)
				assert.Equal(t, tc.err, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.user, got)
			}
		})
	}
}

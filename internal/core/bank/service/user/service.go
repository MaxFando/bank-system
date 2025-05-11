//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mock_${GOFILE}.go -package=${GOPACKAGE}
package user

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"log/slog"
)

type Repository interface {
	UpdateProfile(ctx context.Context, user *entity.User) (*entity.User, error)
}

type Service struct {
	repo   Repository
	logger *slog.Logger
}

func NewService(logger *slog.Logger, repo Repository) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
	}
}

func (s *Service) UpdateProfile(ctx context.Context, user *entity.User) (*entity.User, error) {
	user, err := s.repo.UpdateProfile(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return user, nil
}

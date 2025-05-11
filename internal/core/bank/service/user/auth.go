//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mock_${GOFILE}.go -package=${GOPACKAGE}
package user

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/auth"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

// AuthRepository представляет интерфейс для управления авторизацией пользователей.
// Используется для регистрации новых пользователей и поиска по email.
type AuthRepository interface {
	Register(ctx context.Context, email, password string) (int32, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

// AuthService предоставляет функциональность для регистрации и аутентификации пользователей.
// Включает взаимодействие с хранилищем пользователей (AuthRepository) и логгером.
type AuthService struct {
	repo   AuthRepository
	logger *slog.Logger
}

// NewAuthService создает новый экземпляр AuthService с заданным репозиторием и логгером.
func NewAuthService(logger *slog.Logger, repo AuthRepository) *AuthService {
	return &AuthService{
		repo:   repo,
		logger: logger,
	}
}

// Register регистрирует нового пользователя с переданным email и паролем и возвращает его уникальный идентификатор.
func (s *AuthService) Register(ctx context.Context, email, password string) (int32, error) {
	id, err := s.repo.Register(ctx, email, password)
	if err != nil {
		return 0, fmt.Errorf("failed to register user: %w", err)
	}

	return id, nil
}

// Login выполняет аутентификацию пользователя по email и паролю и возвращает сгенерированный JWT токен.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to find user by email: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	token, err := auth.GenerateJWTToken(fmt.Sprint(user.ID))
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return token, nil
}

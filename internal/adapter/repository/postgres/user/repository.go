package user

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Repository struct {
	db sqlext.DB
}

func NewRepository(db sqlext.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Register(ctx context.Context, email, password string) (int32, error) {
	query := `
		INSERT INTO main.users (email, password_hash)
		VALUES (:email, :password)
		RETURNING id
	`

	// Хеширование пароля с помощью bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	args := map[string]interface{}{
		"email":         email,
		"password_hash": hashedPassword,
	}

	var id int32
	err = r.db.Get(ctx, &id, query, args)
	if err != nil {
		return 0, fmt.Errorf("failed to register user: %w", err)
	}

	return id, nil
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, date_of_birth
		FROM main.users
		WHERE email = :email
	`

	args := map[string]interface{}{
		"email": email,
	}

	var user entity.User
	err := r.db.Get(ctx, &user, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return &user, nil
}

// UpdateProfile обновляет профиль пользователя в базе данных.
func (r *Repository) UpdateProfile(ctx context.Context, user *entity.User) (*entity.User, error) {
	query := `
		UPDATE main.users
		SET first_name    = :first_name,
			last_name     = :last_name,
			date_of_birth = :date_of_birth
		WHERE id = :id
		RETURNING id, first_name, last_name, date_of_birth
	`
	query, args, err := r.db.BindNamed(query, user)
	if err != nil {
		return nil, fmt.Errorf("failed to bind named parameters: %w", err)
	}

	type row struct {
		ID          int32     `db:"id"`
		FirstName   string    `db:"first_name"`
		LastName    string    `db:"last_name"`
		DateOfBirth time.Time `db:"date_of_birth"`
	}

	var updatedUser row
	err = r.db.Get(ctx, &updatedUser, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	user.ID = updatedUser.ID
	user.FirstName = updatedUser.FirstName
	user.LastName = updatedUser.LastName
	user.DateOfBirth = updatedUser.DateOfBirth

	return user, nil
}

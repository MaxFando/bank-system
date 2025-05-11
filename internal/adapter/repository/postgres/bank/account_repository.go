package bank

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"
)

type AccountRepository struct {
	db sqlext.DB
}

func NewAccountRepository(db sqlext.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Save(ctx context.Context, account *entity.Account) (*entity.Account, error) {
	query := `
		INSERT INTO main.accounts (user_id, account_number, balance, account_type)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, account_number, balance, account_type;
	`

	err := r.db.Get(ctx, account, query, account.UserID, account.AccountNumber, account.Balance, account.AccountType)
	if err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) FindByID(ctx context.Context, id int32) (*entity.Account, error) {
	query := `
		SELECT id, user_id, account_number, balance, account_type
		FROM main.accounts
		WHERE id = $1;
	`

	account := &entity.Account{}
	err := r.db.Get(ctx, account, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find account by ID: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) GetAccountByUserID(ctx context.Context, userID int32) (*entity.Account, error) {
	query := `
		SELECT id, user_id, account_number, balance, account_type
		FROM main.accounts
		WHERE user_id = $1;
	`

	account := &entity.Account{}
	err := r.db.Get(ctx, account, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find account by user ID: %w", err)
	}

	return account, nil
}

func (r *AccountRepository) WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	return r.db.WithTx(ctx, fn, opts...)
}

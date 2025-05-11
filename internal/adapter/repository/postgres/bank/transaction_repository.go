package bank

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext"
	"github.com/MaxFando/bank-system/pkg/sqlext/transaction"
)

type CardTransactionRepository struct {
	db sqlext.DB
}

func NewCardTransactionRepository(db sqlext.DB) *CardTransactionRepository {
	return &CardTransactionRepository{
		db: db,
	}
}

func (c CardTransactionRepository) Transfer(ctx context.Context, cardID int32, amount float64) (int32, error) {
	query := `
		INSERT INTO main.card_transactions (card_id, amount, transaction_type)
		VALUES ($1, $2, 'transfer')
		RETURNING id;
	`

	var id int32
	err := c.db.Get(ctx, &id, query, cardID, amount)
	if err != nil {
		return 0, fmt.Errorf("failed to save card transaction: %w", err)
	}

	return id, nil
}

func (c CardTransactionRepository) Withdraw(ctx context.Context, cardID int32, amount float64) (int32, error) {
	query := `
		INSERT INTO main.card_transactions (card_id, amount, transaction_type)
		VALUES ($1, $2, 'withdraw')
		RETURNING id;
	`

	var id int32
	err := c.db.Get(ctx, &id, query, cardID, amount)
	if err != nil {
		return 0, fmt.Errorf("failed to save card transaction: %w", err)
	}

	return id, nil
}

func (c CardTransactionRepository) Deposit(ctx context.Context, cardID int32, amount float64) (int32, error) {
	query := `
		INSERT INTO main.card_transactions (card_id, amount, transaction_type)
		VALUES ($1, $2, 'deposit')
		RETURNING id;
	`

	var id int32
	err := c.db.Get(ctx, &id, query, cardID, amount)
	if err != nil {
		return 0, fmt.Errorf("failed to save card transaction: %w", err)
	}

	return id, nil
}

func (c CardTransactionRepository) FindByID(ctx context.Context, id int32) (*entity.CardTransaction, error) {
	query := `
		SELECT id, card_id, amount, transaction_type
		FROM main.card_transactions
		WHERE id = $1;
	`

	cardTransaction := &entity.CardTransaction{}
	err := c.db.Get(ctx, cardTransaction, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find card transaction by ID: %w", err)
	}

	return cardTransaction, nil
}

func (c CardTransactionRepository) WithTx(ctx context.Context, fn transaction.AtomicFn, opts ...transaction.TxOption) error {
	return c.db.WithTx(ctx, fn, opts...)
}

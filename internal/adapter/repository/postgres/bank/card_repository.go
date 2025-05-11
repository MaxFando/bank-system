package bank

import (
	"context"
	"fmt"
	"github.com/MaxFando/bank-system/internal/core/bank/entity"
	"github.com/MaxFando/bank-system/pkg/sqlext"
)

type CardRepository struct {
	db sqlext.DB
}

func NewCardRepository(db sqlext.DB) *CardRepository {
	return &CardRepository{
		db: db,
	}
}

func (c CardRepository) Save(ctx context.Context, card *entity.Card) (*entity.Card, error) {
	query := `INSERT INTO main.cards (account_id, encrypted_data, hmac) VALUES ($1, $2, $3) RETURNING id`

	var id int32
	err := c.db.Get(ctx, &id, query, card.AccountID, card.EncryptedData, card.HMAC)
	if err != nil {
		return nil, fmt.Errorf("failed to save card: %w", err)
	}

	card.ID = id

	return card, nil
}

func (c CardRepository) FindByID(ctx context.Context, id int32) (*entity.Card, error) {
	query := `SELECT id, account_id, encrypted_data, hmac FROM main.cards WHERE id = $1`

	card := &entity.Card{}
	err := c.db.Get(ctx, card, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find card by ID: %w", err)
	}

	return card, nil
}

func (c CardRepository) FindByAccountID(ctx context.Context, accountID int32) ([]entity.Card, error) {
	query := `SELECT id, account_id, encrypted_data, hmac FROM main.cards WHERE account_id = $1`

	var cards []entity.Card
	err := c.db.Select(ctx, &cards, query, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to find cards by account ID: %w", err)
	}

	return cards, nil
}
